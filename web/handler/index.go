package handler

import (
	"net/http"

	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/labstack/echo/v4"
)

/* Public */

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", echo.Map{
		"apikey":  config.Config.Server.ApiKey,
		"baseurl": config.Config.Server.BaseURL,
	})
}
