package apis

import (
	"github.com/labstack/echo"
	"net/http"
)

/* Structs */

type Response struct {
	Error   bool   `json:"error" xml:"error"`
	Message string `json:"message" xml:"message"`
}

/* Public */

func Test(c echo.Context) error {
	return c.JSON(http.StatusOK, &Response{Error: false, Message: "OK"})
}
