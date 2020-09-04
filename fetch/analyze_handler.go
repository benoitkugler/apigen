package fetch

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/benoitkugler/apigen/gents"
)

// Look for Bind(), QueryParam() and JSON() method calls
// pkg is the pacakge of the method
func analyzeHandler(body []ast.Stmt, pkg *types.Package) gents.Contrat {
	var out gents.Contrat
	for _, stmt := range body {
		switch stmt := stmt.(type) {
		case *ast.ReturnStmt:
			if len(stmt.Results) != 1 { // should not happend : the method return error
				continue
			}
			if call, ok := stmt.Results[0].(*ast.CallExpr); ok {
				if method, ok := call.Fun.(*ast.SelectorExpr); ok {
					if method.Sel.Name == "JSON" || method.Sel.Name == "JSONPretty" {
						if len(call.Args) >= 2 {
							output := call.Args[1] // c.JSON(200, output)
							if output, ok := output.(*ast.Ident); ok {
								out.Return = resolveLocalType(output, pkg)
							}
						}
					}
				}
			}

		case *ast.AssignStmt:
			for _, rh := range stmt.Rhs {
				out.Input = parseBindCall(rh, pkg)
				if queryParam := parseCallWithString(rh, "QueryParam"); queryParam != "" {
					out.QueryParams = append(out.QueryParams, queryParam)
				}
				if formValue := parseCallWithString(rh, "FormValue"); formValue != "" {
					out.Form.Values = append(out.Form.Values, formValue)
				}
				if formFile := parseCallWithString(rh, "FormFile"); formFile != "" {
					out.Form.File = formFile
				}
			}
		case *ast.IfStmt:
			if asign, ok := stmt.Init.(*ast.AssignStmt); ok {
				for _, rh := range asign.Rhs {
					out.Input = parseBindCall(rh, pkg)
					if formFile := parseCallWithString(rh, "FormFile"); formFile != "" {
						out.Form.File = formFile
					}
				}
			}

		}
	}
	return out
}

func parseBindCall(expr ast.Expr, pkg *types.Package) types.Type {
	if call, ok := expr.(*ast.CallExpr); ok {
		if caller, ok := call.Fun.(*ast.SelectorExpr); ok {
			if caller.Sel.Name == "Bind" && len(call.Args) == 1 { // "c.Bind(in)"
				switch arg := call.Args[0].(type) {
				case *ast.Ident: // c.Bind(pointer)
					return resolveLocalType(arg, pkg)
				case *ast.UnaryExpr: // c.Bind(&value)
					if ident, ok := arg.X.(*ast.Ident); arg.Op == token.AND && ok {
						return resolveLocalType(ident, pkg)
					}
				}
			}
		}
	}
	return nil
}

func parseCallWithString(expr ast.Expr, methodName string) string {
	if call, ok := expr.(*ast.CallExpr); ok {
		if caller, ok := call.Fun.(*ast.SelectorExpr); ok {
			if caller.Sel.Name == methodName && len(call.Args) == 1 { // "c.<methodName>(<string>)"
				if lit, ok := call.Args[0].(*ast.BasicLit); ok {
					return stringLitteral(lit)
				}
			}
		}
	}
	return ""
}

func resolveLocalType(arg *ast.Ident, pkg *types.Package) types.Type {
	localScope := pkg.Scope().Innermost(arg.Pos())
	obj := localScope.Lookup(arg.Name)
	for obj == nil {
		localScope = localScope.Parent()
		obj = localScope.Lookup(arg.Name)
	}
	return obj.Type()
}
