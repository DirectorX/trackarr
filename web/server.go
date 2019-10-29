package web

import (
	"fmt"
	"net/http"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/runtime"
	"github.com/l3uddz/trackarr/web/apis"
	"github.com/l3uddz/trackarr/web/handler"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/echo-template/supports/gorice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pascaldekloe/latest"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
)

/* Structs */

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

/* Vars */
var (
	log        = logger.GetLogger("web")
	logEmitter latest.Broadcast
)

/* Public */

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Listen(configuration *config.Configuration, logLevel int) {
	/* init echo */
	e := echo.New()
	if logLevel > 1 {
		// log to stdout
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	// setup validator
	e.Validator = &CustomValidator{
		validator: validator.New(&validator.Config{
			TagName:      "validate",
			FieldNameTag: "validate",
		}),
	}

	// setup template renderer
	e.Renderer = gorice.New(rice.MustFindBox("views"))

	// setup static file server
	staticBox := rice.MustFindBox("static")
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))

	// setup groups
	gui := e.Group("")
	if configuration.Server.User != "" && configuration.Server.Pass != "" {
		// user and pass were defined, use basic auth middleware
		gui.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == configuration.Server.User && password == configuration.Server.Pass {
				return true, nil
			}
			return false, nil
		}))
	}

	api := e.Group("/api", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:apikey",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == configuration.Server.ApiKey, nil
		},
	}))

	// - add api routes
	api.GET("/torrent", apis.Torrent)
	api.GET("/releases", apis.Releases)
	api.GET("/irc/status", apis.IrcStatus)

	/* init frontend routes */
	// static
	gui.GET("/static/*", echo.WrapHandler(staticFileServer))

	// index
	gui.GET("/", handler.Index)

	// logs
	gui.GET("/logs", handler.Logs)
	gui.GET("/logs/ws", WebsocketLogHandler)

	// close broadcaster
	defer logEmitter.UnsubscribeAll()

	// setup log hook and emitter
	logrus.AddHook(&WebsocketLogHook{})

	/* start echo server */
	runtime.Web = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", configuration.Server.Host, configuration.Server.Port),
		Handler: e,
	}

	go func() {
		log.Infof("Listening on %s:%d", configuration.Server.Host, configuration.Server.Port)

		if err := runtime.Web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatalf("Failed listening on %s:%d", configuration.Server.Host, configuration.Server.Port)
		}
	}()
}
