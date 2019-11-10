package ws

import "github.com/labstack/echo"

func HandlerFunc(context echo.Context) error {
	return m.HandleRequest(context.Response().Writer, context.Request())
}
