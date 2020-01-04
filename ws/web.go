package ws

import (
	"github.com/labstack/echo/v4"
)

func HandlerFunc(context echo.Context) error {
	return m.HandleRequest(context.Response().Writer, context.Request())
}
