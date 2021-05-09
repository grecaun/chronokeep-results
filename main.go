package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"chronokeep/results/database"
	"chronokeep/results/handlers"
	"chronokeep/results/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := util.GetConfig()
	if err != nil {
		log.Fatal("Failed to get configuration. ", err)
	}
	e := echo.New()
	e.Debug = config.Development

	err = database.Setup(config)
	if err != nil {
		log.Fatal("Unable to setup database. ", err)
	}

	// Set up Recover and Logger middleware.
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:  "${method} | ${status} | ${uri} -> ${latency_human}\n",
		Skipper: healthEndpointSkipper,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
	}))

	handler := handlers.Handler{}
	handler.Bind(e.Group(""))
	r := e.Group("")
	r.Use(middleware.JWT([]byte(config.SecretKey)))
	handler.BindRestricted(r)

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	if config.Development {
		e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Handler: func(c echo.Context, req []byte, res []byte) {
				log.Info("Request Log: ", string(req))
				log.Info("Request Header: ", c.Request().Header)
				log.Info("Response Log: ", string(res))
			},
			Skipper: healthEndpointSkipper,
		}))
	}
	if !config.AutoTLS {
		log.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
	} else {
		// Set up TLS with auto certificate if not a debug environment.
		// e.AutoTLSManager.HostPolicy = autocert.HostWhiteList("<DOMAIN>")
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Pre(middleware.HTTPSRedirect())
		log.Fatal(e.StartAutoTLS(":" + strconv.Itoa(config.Port)))
	}
}

func healthEndpointSkipper(c echo.Context) bool {
	return strings.HasPrefix(c.Path(), "/health")
}

func init() {
	config, err := util.GetConfig()
	if err != nil {
		log.Fatal("Failed to get configuration. ", err)
	}
	if config.Development {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}
}
