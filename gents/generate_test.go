package gents

import (
	"fmt"
	"go/types"
	"net/http"
	"testing"

	"github.com/benoitkugler/structgen/tstypes"
)

func TestGenerate(t *testing.T) {
	apis := Service{
		{
			Url: "/samlskm/", Method: http.MethodPost, Contrat: Contrat{
				Input:       TypeNoId{Type: types.NewArray(types.Typ[types.Byte], 5), NoId: true},
				HandlerName: "M1",
				Return:      types.NewSlice(types.Typ[types.Int]),
			},
		},
		{
			Url: "/samlskm/:param1", Method: http.MethodGet, Contrat: Contrat{
				HandlerName: "M2",
				QueryParams: []TypedParam{
					{Name: "arg1", Type: tstypes.TsString},
					{Name: "arg2", Type: tstypes.TsNumber},
					{Name: "arg3", Type: tstypes.TsBoolean},
				},
				Return: types.NewSlice(types.Typ[types.Int]),
			},
		},
	}
	fmt.Println(apis.Render(nil))
}
