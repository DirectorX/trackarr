package web

import (
	"fmt"
	"github.com/l3uddz/trackarr/ws"
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
	log = logger.GetLogger("web")
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
	e.Use(middleware.CORS())

	// setup validator
	e.Validator = &CustomValidator{
		validator: validator.New(&validator.Config{
			TagName:      "validate",
			FieldNameTag: "validate",
		}),
	}

	// setup template renderer
	e.Renderer = gorice.New(rice.MustFindBox("trackarr-ui/dist"))

	// setup websocket server
	if err := ws.Init(); err != nil {
		log.WithError(err).Fatal("Failed initializing websocket server")
	}

	// setup static file server
	staticBox := rice.MustFindBox("trackarr-ui/dist/static")
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
	api.Any("/ws", ws.HandlerFunc)
	api.GET("/torrent", apis.Torrent)
	api.GET("/releases", apis.Releases)
	api.GET("/irc/status", apis.IrcStatus)

	/* init frontend routes */
	// static
	gui.GET("/static/*", echo.WrapHandler(staticFileServer))

	// ui
	gui.Any("/*", handler.Index)

	// setup log hook
	if err := runtime.Loghook.Start(); err != nil {
		log.WithError(err).Error("Failed starting loghook")
	} else {
		log.Info("Started loghook")
	}

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
