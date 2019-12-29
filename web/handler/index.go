package handler

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"github.com/labstack/echo"
	"net/http"
)

/* Public */

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", echo.Map{
		"apikey": config.Config.Server.ApiKey,
	})
}
