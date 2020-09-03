package main

import (
	"github.com/benoitkugler/apigen/fetch/test/inner"
	"github.com/labstack/echo"
)

const route = "/const_url_from_package/"

type controller struct {
}

func (controller) M(c echo.Context) error {
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

func routes(e *echo.Echo, ct controller, ct2 inner.Controller) {
	e.GET(route, handler)
	const routeFunc = "const_local_url"
	e.GET(routeFunc, ct.M)
	e.POST(inner.Url, ct2.M)
	e.POST(inner.Url+"endpoint", ct2.M)
	e.POST("host"+inner.Url, ct2.M)
	e.POST("host"+"endpoint", ct2.M)
	e.POST("/string_litteral", ct2.M)
	e.PUT("/with_param/:param", ct2.M)
	e.DELETE("/special_param_value/:class/route", ct2.M)
	e.DELETE("/special_param_value/:default/route", ct2.M)
}
