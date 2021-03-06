package gents

import (
	"fmt"
	"go/types"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/tstypes"
)

type TypedParam struct {
	Name string
	Type tstypes.TsType
}

// return arg: String(params[arg])
func (t TypedParam) asObjectKey() string {
	out := fmt.Sprintf("%q: ", t.Name)
	switch t.Type {
	case tstypes.TsNumber:
		out += fmt.Sprintf("String(params[%q])", t.Name) // stringify
	case tstypes.TsBoolean:
		out += fmt.Sprintf("params[%q] ? 'ok' : ''", t.Name) // stringify
	default:
		out += fmt.Sprintf("params[%q]", t.Name) // no converter
	}
	return out
}

// TypeNoId represents a type that might omit the "id" field
type TypeNoId struct {
	Type types.Type
	NoId bool
}

func (t TypeNoId) render() string {
	baseType := goToTs(t.Type).Render()
	if t.NoId {
		return "New<" + baseType + ">"
	}
	return baseType
}

type Contrat struct {
	HandlerName string
	Input       TypeNoId
	QueryParams []TypedParam
	Form        Form
	Return      types.Type
}

type API struct {
	Url     string
	Method  string
	Contrat Contrat
}

type Form struct {
	Values []string
	File   string // empty means no file
}

func (f Form) IsZero() bool {
	return f.File == "" && len(f.Values) == 0
}

// type string
func (f Form) typedValues() []TypedParam {
	out := make([]TypedParam, len(f.Values))
	for i, v := range f.Values {
		out[i] = TypedParam{Name: v, Type: tstypes.TsString}
	}
	return out
}

func (a API) withBody() bool {
	return a.Method == http.MethodPost || a.Method == http.MethodPut
}

func (a API) withFormData() bool {
	return !a.Contrat.Form.IsZero()
}

func (a API) methodLower() string {
	return strings.ToLower(a.Method)
}

func paramsType(params []TypedParam) string {
	tmp := make([]string, len(params))
	for i, s := range params {
		tmp[i] = fmt.Sprintf("%q: %s", s.Name, s.Type.Render()) // quote for names like "id-1"
	}
	return "{" + strings.Join(tmp, ", ") + "}"
}

func (a API) funcArgsName() string {
	if a.withBody() {
		if a.withFormData() { // form data mode
			if fi := a.Contrat.Form.File; fi != "" {
				return "params, file"
			}
		}
	} else {
		// params as query params
		if len(a.Contrat.QueryParams) == 0 {
			return ""
		}
	}
	return "params"
}

func goToTs(typ types.Type) tstypes.TsType {
	return tstypes.NewHandler(nil).AnalyseType(loader.NewDeclarations(), typ)
}

func (a API) typeIn() string {
	if a.withBody() {
		if a.withFormData() { // form data mode
			params := "params: " + paramsType(a.Contrat.Form.typedValues())
			if fi := a.Contrat.Form.File; fi != "" {
				params += ", file: File"
			}
			return params
		} else { // JSON mode
			return "params: " + a.Contrat.Input.render()
		}
	}
	// params as query params
	if len(a.Contrat.QueryParams) == 0 {
		return ""
	}
	return "params: " + paramsType(a.Contrat.QueryParams)
}

// use a named package
func (a API) typeOut() string {
	return goToTs(a.Contrat.Return).Render()
}

var rePlaceholder = regexp.MustCompile(`:([^/"']+)`)

const templateFuncReplace = `(%s) => %s%s` // path ,  .replace(placeholder, args[0]) ...

// returns the names of the params in url, in two versions:
// the original ones (with :) and ts compatible ones
func (a API) parseUrlParams() ([]string, []string) {
	pls := rePlaceholder.FindAllString(a.Url, -1)
	tsCompatible := make([]string, len(pls))
	for i, pl := range pls {
		argname := pl[1:]
		if argname == "default" || argname == "class" { // js keywords
			argname += "_"
		}
		tsCompatible[i] = argname
	}
	return pls, tsCompatible
}

func (a API) fullUrl() string {
	params, names := a.parseUrlParams()
	if len(params) > 0 {
		// the url has arguments
		var calls string
		for i, placeholder := range params {
			calls += fmt.Sprintf(".replace('%s', this.urlParams.%s)", placeholder, names[i])
		}
		return fmt.Sprintf("this.baseUrl + %q%s", a.Url, calls)
	}
	return fmt.Sprintf("this.baseUrl + %q", a.Url) // basic url
}

func (c Contrat) convertTypedQueryParams() string {
	chunks := make([]string, len(c.QueryParams))
	for i, param := range c.QueryParams {
		chunks[i] = param.asObjectKey()
	}
	return "{ " + strings.Join(chunks, ", ") + " }"
}

func (a API) generateCall() string {
	var template string
	if a.withBody() {
		if a.withFormData() { // add the creation of FormData
			template += "const formData = new FormData()\n"
			if fi := a.Contrat.Form.File; fi != "" {
				template += fmt.Sprintf("formData.append(%q, file, file.name)\n", fi)
			}
			for _, param := range a.Contrat.Form.Values {
				template += fmt.Sprintf("formData.append(%q, params[%q])\n", param, param)
			}
			template += "const rep:AxiosResponse<%s> = await Axios.%s(fullUrl, formData, { headers : this.getHeaders() })"
		} else {
			template = "const rep:AxiosResponse<%s> = await Axios.%s(fullUrl, params, { headers : this.getHeaders() })"
		}
	} else {
		callParams := ", { headers: this.getHeaders() }"
		if len(a.Contrat.QueryParams) != 0 {
			callParams = fmt.Sprintf(", { params: %s, headers : this.getHeaders() }", a.Contrat.convertTypedQueryParams())
		}
		template = "const rep:AxiosResponse<%s> = await Axios.%s(fullUrl" + callParams + ")"
	}
	return fmt.Sprintf(template, a.typeOut(), a.methodLower())
}

func (a API) generateMethod() string {
	const template = `
	protected async raw%s(%s) {
		const fullUrl = %s;
		%s;
		return rep.data;
	}
	
	/** %s wraps raw%s and handles the error */
	async %s(%s) {
		this.startRequest();
		try {
			const out = await this.raw%s(%s);
			this.onSuccess%s(out);
			return out
		} catch (error) {
			this.handleError(error);
		}
	}

	protected abstract onSuccess%s(data: %s): void 
	`
	fnName := a.Contrat.HandlerName
	return fmt.Sprintf(template,
		fnName, a.typeIn(), a.fullUrl(), a.generateCall(), fnName, fnName, fnName, a.typeIn(),
		fnName, a.funcArgsName(), fnName, fnName, a.typeOut())
}

type Service []API

// aggregate the different url params
func (s Service) urlParamsType() string {
	all := map[string]bool{}
	for _, api := range s {
		_, params := api.parseUrlParams()
		for _, param := range params {
			all[param] = true
		}
	}
	sorted := make(sort.StringSlice, 0, len(all))
	for param := range all {
		sorted = append(sorted, param)
	}
	sort.Sort(sorted)
	for i, param := range sorted {
		sorted[i] = param + ": string"
	}
	return "{" + strings.Join(sorted, ", ") + "}"
}

type tsNewDeclaration struct{}

func (tsNewDeclaration) Id() string { return "__ts_new_declaration" }
func (tsNewDeclaration) Render() string {
	return `export type New<T extends { id: number }> = Omit<T, "id"> & Partial<Pick<T, "id">>;`
}

func (s Service) renderTypes(enum enums.EnumTable) string {
	decl := loader.NewDeclarations()
	handler := tstypes.NewHandler(enum)
	for _, api := range s { // write top-level decl
		handler.AnalyseType(decl, api.Contrat.Input.Type)
		if api.Contrat.Input.NoId {
			decl.Add(tsNewDeclaration{})
		}
		handler.AnalyseType(decl, api.Contrat.Return)
	}
	return decl.ToString()
}

func (s Service) Render(enum enums.EnumTable) string {
	apiCalls := make([]string, len(s))
	for i, api := range s {
		apiCalls[i] = api.generateMethod()
	}
	return fmt.Sprintf(`
	// Code generated by apigen. DO NOT EDIT
	
	import Axios, { AxiosResponse } from "axios";

	%s

	export abstract class AbstractAPI {
		constructor(protected baseUrl: string, protected authToken: string, protected urlParams: %s) {}

		abstract handleError(error: any): void

		abstract startRequest(): void

		getHeaders() {
			return { Authorization: "Bearer " + this.authToken }
		}

		%s
	}`, s.renderTypes(enum), s.urlParamsType(), strings.Join(apiCalls, "\n"))
}
