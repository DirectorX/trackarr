package handler

import (
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/labstack/echo"
	"net/http"
)

/* Public */

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", echo.Map{
		"title":    "Home",
		"page":     "index",
		"latest":   models.GetLatestPushedReleases(database.DB, 20),
		"approved": models.GetLatestApprovedReleases(database.DB, 20),
	})
}
