package gents

import (
	"fmt"
	"go/types"
	"net/http"
)

type Contrat struct {
	HandlerName string
	Input       types.Type
	QueryParams []string
	Return      types.Type
}

type API struct {
	Url     string
	Method  string
	Contrat Contrat
}

func (a API) generateCall() string {
	if a.Method == http.MethodGet || a.Method == http.MethodDelete {
		return fmt.Sprintf("const rep:<AxiosResponse<%s>> = await Axios.%s(%s, {query: params})", typeOut, methodLower, url)
	} else {
		return fmt.Sprintf("const rep:<AxiosResponse<%s>> = await Axios.%s(%s, params)", typeOut, methodLower, url)
	}
}

func (a API) generateMethod() {

}
