package gents

import (
	"fmt"
	"go/types"
	"net/http"
	"testing"
)

func TestGenerate(t *testing.T) {
	apis := Service{
		{
			Url: "/samlskm/", Method: http.MethodGet, Contrat: Contrat{
				QueryParams: []string{"arg1", "arg2"},
				Return:      types.NewSlice(types.Typ[types.Int]),
			},
		},
		{
			Url: "/samlskm/:param1", Method: http.MethodGet, Contrat: Contrat{
				QueryParams: []string{"arg1", "arg2"},
				Return:      types.NewSlice(types.Typ[types.Int]),
			},
		},
	}
	fmt.Println(apis.Render())
}
