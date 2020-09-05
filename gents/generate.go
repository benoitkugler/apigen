package gents

import (
	"fmt"
	"go/types"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/benoitkugler/structgen/tstypes"
)

type Contrat struct {
	HandlerName string
	Input       types.Type
	QueryParams []string
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

func (a API) withBody() bool {
	return a.Method == http.MethodPost || a.Method == http.MethodPut
}

func (a API) withFormData() bool {
	return !a.Contrat.Form.IsZero()
}

func (a API) methodLower() string {
	return strings.ToLower(a.Method)
}

func paramsType(params []string) string {
	tmp := make([]string, len(params))
	for i, s := range params {
		tmp[i] = fmt.Sprintf("%q: string", s) // quote for names like "id-1"
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

func (a API) typeIn() string {
	if a.withBody() {
		if a.withFormData() { // form data mode
			params := "params: " + paramsType(a.Contrat.Form.Values)
			if fi := a.Contrat.Form.File; fi != "" {
				params += ", file: File"
			}
			return params
		} else { // JSON mode
			return "params: " + tstypes.GoToTs(a.Contrat.Input).Render()
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
	ts := tstypes.GoToTs(a.Contrat.Return)
	if _, isNamed := ts.(tstypes.TsNamedType); isNamed {
		return "types." + ts.Render()
	}
	return ts.Render()
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
		var queryParams string
		if len(a.Contrat.QueryParams) != 0 {
			queryParams = ", { params: params, headers : this.getHeaders() }"
		}
		template = "const rep:AxiosResponse<%s> = await Axios.%s(fullUrl" + queryParams + ")"
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
	
	async %s(%s) {
		this.startRequest();
		try {
			const out = await this.raw%s(%s);
			this.onSuccess%s(out);
		} catch (error) {
			this.handleError(error);
		}
	}

	protected abstract onSuccess%s(data: %s): void 
	`
	fnName := a.Contrat.HandlerName
	return fmt.Sprintf(template,
		fnName, a.typeIn(), a.fullUrl(), a.generateCall(), fnName, a.typeIn(),
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

func (s Service) Render() string {
	apiCalls := make([]string, len(s))
	for i, api := range s {
		apiCalls[i] = api.generateMethod()
	}
	return fmt.Sprintf(`
	// Code generated by apigen. DO NOT EDIT
	
	import Axios, { AxiosResponse } from "axios";

	export abstract class AbstractAPI {
		constructor(protected baseUrl: string, protected authToken: string, protected urlParams: %s) {}

		abstract handleError(error: any): void

		abstract startRequest(): void

		getHeaders() {
			return { Authorization: "Bearer " + this.authToken }
		}

		%s
	}`, s.urlParamsType(), strings.Join(apiCalls, "\n"))
}
