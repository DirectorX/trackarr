package web

import (
	"fmt"
	"net/http"
	"path"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/runtime"
	"gitlab.com/cloudb0x/trackarr/web/apis"
	"gitlab.com/cloudb0x/trackarr/web/handler"
	"gitlab.com/cloudb0x/trackarr/ws"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/foolin/goview/supports/gorice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func Listen(cfg *config.Configuration, logLevel int) {
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
	ev := echoview.Default()
	ev.SetFileHandler(gorice.FileHandler(rice.MustFindBox("trackarr-ui/dist")))
	e.Renderer = ev

	// setup websocket server
	if err := ws.Init(); err != nil {
		log.WithError(err).Fatal("Failed initializing websocket server")
	}

	// setup static file server
	staticFileServer := http.StripPrefix(
		path.Join(cfg.Server.BaseURL, "/static/"),
		http.FileServer(
			rice.MustFindBox("trackarr-ui/dist/static").HTTPBox(),
		),
	)

	// UI
	var gui *echo.Group
	if cfg.Server.BaseURL == "/" {
		gui = e.Group("")
	} else {
		gui = e.Group(cfg.Server.BaseURL)
	}
	// Basic auth
	if cfg.Server.User != "" && cfg.Server.Pass != "" {
		// user and pass were defined, use basic auth middleware
		gui.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == cfg.Server.User && password == cfg.Server.Pass {
				return true, nil
			}
			return false, nil
		}))
	}

	// Redirect to base URL when not `/`
	if cfg.Server.BaseURL != "/" {
		// Root to base URL
		e.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, cfg.Server.BaseURL+"/")
		})
		// Add `/` to base URL
		gui.GET("", handler.Index, middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
			RedirectCode: http.StatusFound,
		}))
	}

	// - UI routes
	gui.Any("/*", handler.Index)
	gui.GET("/static/*", echo.WrapHandler(staticFileServer))

	// API
	api := e.Group(path.Join(cfg.Server.BaseURL, "/api"),
		middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:apikey",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == cfg.Server.ApiKey, nil
			},
		}),
	)

	// - API routes
	api.Any("/ws", ws.HandlerFunc)
	api.GET("/torrent", apis.Torrent)
	api.GET("/releases", apis.Releases)
	api.GET("/irc/status", apis.IrcStatus)
	api.GET("/update/status", apis.UpdateStatus)

	// setup log hook
	if err := runtime.Loghook.Start(); err != nil {
		log.WithError(err).Error("Failed starting loghook")
	} else {
		log.Info("Started loghook")
	}

	/* start echo server */
	runtime.Web = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: e,
	}

	log.Infof("Listening on %s:%d%s", cfg.Server.Host, cfg.Server.Port, cfg.Server.BaseURL)

	go func() {
		if err := runtime.Web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatalf("Failed listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		}
	}()
}
