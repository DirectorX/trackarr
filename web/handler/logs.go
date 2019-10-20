package handler

import (
	"github.com/labstack/echo"
	"net/http"
)

/* Public */

func Logs(c echo.Context) error {
	return c.Render(http.StatusOK, "logs", echo.Map{
		"title": "Logs",
		"page":  "logs",
	})
}

