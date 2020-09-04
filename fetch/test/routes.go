package main

import (
	"fmt"

	"github.com/benoitkugler/apigen/fetch/test/inner"
	"github.com/labstack/echo"
)

const route = "/const_url_from_package/"

type controller struct {
}

func (controller) handle1(c echo.Context) error {
	var (
		in  int
		out string
	)
	if err := c.Bind(&in); err != nil {
		return err
	}
	return c.JSON(200, out)
}

func handler(echo.Context) error { return nil }

func (controller) handler2(echo.Context) error { return nil }
func (controller) handler3(echo.Context) error { return nil }
func (controller) handler4(echo.Context) error { return nil }
func (controller) handler5(echo.Context) error { return nil }
func (controller) handler6(echo.Context) error { return nil }
func (controller) handler7(echo.Context) error { return nil }
func (controller) handler8(c echo.Context) error {
	id1, id2 := c.QueryParam("query_param1"), c.QueryParam("query_param2")
	fmt.Println(id1, id2)
	var code uint
	return c.JSON(200, code)
}

func routes(e *echo.Echo, ct controller, ct2 inner.Controller) {
	e.GET(route, handler)
	const routeFunc = "const_local_url"
	e.GET(routeFunc, ct.handle1)
	e.POST(inner.Url, ct2.handleExt)
	e.POST(inner.Url+"endpoint", ct.handler2)
	e.POST("host"+inner.Url, ct.handler3)
	e.POST("host"+"endpoint", ct.handler4)
	e.POST("/string_litteral", ct.handler5)
	e.PUT("/with_param/:param", ct.handler6)
	e.DELETE("/special_param_value/:class/route", ct.handler7)
	e.DELETE("/special_param_value/:default/route", ct.handler8)
}
