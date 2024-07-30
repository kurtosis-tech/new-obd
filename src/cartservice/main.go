package main

import (
	"fmt"
	cartservice_server_rest_server "github.com/kurtosis-tech/new-obd/src/cartservice/api/http_rest/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net"
)

const (
	pathToApiGroup         = "/api"
	restAPIPortAddr uint16 = 8090
	restAPIHostIP   string = "0.0.0.0"
)

var (
	defaultCORSOrigins = []string{"*"}
	defaultCORSHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}
)

func main() {
	logrus.Info("Running REST API server...")

	// This is how you set up a basic Echo router
	echoRouter := echo.New()
	echoApiRouter := echoRouter.Group(pathToApiGroup)
	echoApiRouter.Use(middleware.Logger())

	// CORS configuration
	echoApiRouter.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: defaultCORSOrigins,
		AllowHeaders: defaultCORSHeaders,
	}))

	server := NewServer()

	cartservice_server_rest_server.RegisterHandlers(echoApiRouter, cartservice_server_rest_server.NewStrictHandler(server, nil))

	echoRouter.Start(net.JoinHostPort(restAPIHostIP, fmt.Sprint(restAPIPortAddr)))
}
