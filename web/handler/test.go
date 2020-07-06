package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/cloudb0x/trackarr/logger"
	"net/http"
	"strings"
)

/* Public */

type TestRequest struct {
	EventType string
}

type TestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func Test(c echo.Context) error {
	// log
	log := logger.GetLogger("web").WithFields(logrus.Fields{"client": c.RealIP()})

	wr := new(TestRequest)

	if err := c.Bind(wr); err != nil {
		log.WithError(err).Error("Failed parsing test request")
		return c.JSON(http.StatusBadRequest, &TestResponse{
			Status:  "ERROR",
			Message: fmt.Sprintf("error parsing request: %v", err),
		})
	}

	if !strings.EqualFold(wr.EventType, "Test") {
		log.Warnf("Failed validating test request, request: %#v", wr)
		return c.JSON(http.StatusNotAcceptable, &TestResponse{
			Status:  "ERROR",
			Message: fmt.Sprintf("failed validating request event type: %q", wr.EventType),
		})
	}

	log.Infof("Validated test request, request: %+v", wr)
	return c.JSON(http.StatusOK, &TestResponse{
		Status: "OK",
	})
}
