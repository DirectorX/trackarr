package handler

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/labstack/echo"
	"net/http"
)

/* Public */

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", echo.Map{
		"apikey":    config.Config.Server.ApiKey,
	})
}
