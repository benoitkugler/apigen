package fetch

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"regexp"

	"github.com/benoitkugler/apigen/gents"
	"golang.org/x/tools/go/packages"
)

func FetchAPIs(source string) []gents.API {
	pack, file, err := loadSource(source)
	if err != nil {
		log.Fatal(err)
	}
	return parse(pack, file)
}

func isHttpMethod(name string) bool {
	switch name {
	case "GET", "PUT", "POST", "DELETE":
		return true
	default:
		return false
	}
}

func loadSource(sourceFile string) (*packages.Package, *ast.File, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, "file="+sourceFile)
	if err != nil {
		return nil, nil, err
	}
	if len(pkgs) != 1 {
		return nil, nil, fmt.Errorf("only one package expected, got %d", len(pkgs))
	}
	pkg := pkgs[0]
	if L := len(pkg.Syntax); L != 1 {
		return nil, nil, fmt.Errorf("expect one ast.File, got %d", L)
	}
	return pkg, pkg.Syntax[0], nil
}

// look for method calls .GET .POST .PUT .DELETE
// inside all top level functions in `filename`
func parse(pkg *packages.Package, f *ast.File) []gents.API {
	var out []gents.API
	for _, decl := range f.Decls {
		funcStm, ok := decl.(*ast.FuncDecl)
		if !ok || funcStm.Body == nil {
			continue
		}

		for _, stm := range funcStm.Body.List {
			call, ok := stm.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := call.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selector, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			methodName := selector.Sel.Name
			if !isHttpMethod(methodName) || len(callExpr.Args) < 2 {
				// we are looking for .<METHOD>(url, handler)
				continue
			}
			path := parseArgPath(callExpr.Args[0], pkg)
			if path == "" {
				continue
			}
			path = replacePlaceholders(path)
			contrat, err := parseArgHandler(callExpr.Args[1], pkg)
			if err != nil {
				log.Printf("ignoring handler : %s", err)
				continue
			}
			out = append(out, gents.API{Url: path, Method: methodName, Contrat: contrat})
		}
	}
	return out
}

func isImportedPacakge(ident *ast.Ident, pkg *packages.Package) (*packages.Package, bool) {
	for _, imported := range pkg.Imports {
		if imported.Name == ident.Name {
			return imported, true
		}
	}
	return nil, false
}

func resolveStringConst(arg *ast.Ident, pkg *packages.Package) string {
	// start by local scope
	localScope := pkg.Types.Scope().Innermost(arg.Pos())
	var obj types.Object
	if localScope != nil {
		obj = localScope.Lookup(arg.Name)
	}
	if obj == nil { // package scope
		obj = pkg.Types.Scope().Lookup(arg.Name)
	}
	val := obj.(*types.Const).Val()
	if val.Kind() == constant.String {
		return constant.StringVal(val)
	}
	return ""
}

func parseAddStrings(x, y ast.Expr, pkg *packages.Package) string {
	valueX := parseArgPath(x, pkg)
	valueY := parseArgPath(y, pkg)
	return valueX + valueY
}

func stringLitteral(arg *ast.BasicLit) string {
	if arg.Kind == token.STRING { // remove quotes
		return arg.Value[1 : len(arg.Value)-1]
	}
	return ""
}

// we support string litteral or string const
func parseArgPath(arg ast.Expr, pkg *packages.Package) string {
	switch arg := arg.(type) {
	case *ast.Ident:
		if arg.Obj.Kind == ast.Con { // constant of the package
			return resolveStringConst(arg, pkg)
		}
	case *ast.SelectorExpr: // looking for imported constants
		if pkgIdent, ok := arg.X.(*ast.Ident); ok {
			if pkgImported, ok := isImportedPacakge(pkgIdent, pkg); ok {
				return resolveStringConst(arg.Sel, pkgImported)
			}
		}
	case *ast.BinaryExpr:
		if arg.Op == token.ADD {
			return parseAddStrings(arg.X, arg.Y, pkg)
		}
	case *ast.BasicLit:
		if out := stringLitteral(arg); out != "" {
			return out
		}
	}
	log.Printf("Ignoring invalid type for url : %T", arg)
	return ""
}

func resolveMethodReceiver(x *ast.Ident, pkg *packages.Package) *types.Named {
	localScope := pkg.Types.Scope().Innermost(x.Pos())
	obj := localScope.Lookup(x.Name)
	if obj == nil {
		obj = pkg.Types.Scope().Lookup(x.Name)
	}
	if named, ok := obj.Type().(*types.Named); ok {
		return named
	}
	return nil
}

func extractMethodBody(f *ast.File, pos token.Pos) (body []ast.Stmt, err error) {
	for _, decl := range f.Decls {
		funcDecl, isFunc := decl.(*ast.FuncDecl)
		if !isFunc {
			continue
		}
		if funcDecl.Name.NamePos == pos {
			return funcDecl.Body.List, nil
		}
	}
	return nil, errors.New("method not found")
}

// return the file where `fn` is defined
func findMethod(fn *types.Func, rootPkg *packages.Package) (body []ast.Stmt, err error) {
	declFile := rootPkg.Fset.Position(fn.Pos()).Filename

	search := func(pkg *packages.Package) *ast.File {
		for i, file := range pkg.GoFiles {
			if file == declFile {
				return pkg.Syntax[i]
			}
		}
		return nil
	}
	// search in current package
	if f := search(rootPkg); f != nil {
		return extractMethodBody(f, fn.Pos())
	}
	// search into imports
	for _, importedPkg := range rootPkg.Imports {
		if f := search(importedPkg); f != nil {
			return extractMethodBody(f, fn.Pos())
		}
	}
	return nil, errors.New("method not found")
}

func parseArgHandler(arg ast.Expr, pkg *packages.Package) (gents.Contrat, error) {
	if method, ok := arg.(*ast.SelectorExpr); ok {
		if ident, ok := method.X.(*ast.Ident); ok {
			named := resolveMethodReceiver(ident, pkg)
			if named != nil {
				for i := 0; i < named.NumMethods(); i++ {
					fn := named.Method(i)
					if method.Sel.Name == fn.Name() {
						funcBody, err := findMethod(fn, pkg)
						if err != nil {
							return gents.Contrat{}, err
						}
						contrat := analyzeHandler(funcBody, named.Obj().Pkg())
						contrat.HandlerName = fn.Name()
						return contrat, nil
					}
				}
			}
		}
	}
	return gents.Contrat{}, fmt.Errorf("invalid type for handler : %T", arg)
}

var rePlaceholder = regexp.MustCompile(`:([^/"']+)`)

const templateFuncReplace = `(%s) => %s%s` // path ,  .replace(placeholder, args[0]) ...

func replacePlaceholders(endpoint string) string {
	pls := rePlaceholder.FindAllString(endpoint, -1)
	if len(pls) > 0 {
		// the url has arguments
		var args, calls string
		for _, pl := range pls {
			argname := pl[1:]
			if argname == "default" || argname == "class" { // js keywords
				argname += "_"
			}
			args += argname + ":string,"
			calls += fmt.Sprintf(".replace('%s', %s)", pl, argname)
		}
		return fmt.Sprintf(templateFuncReplace, args, endpoint, calls)
	}
	return endpoint // basic url
}
