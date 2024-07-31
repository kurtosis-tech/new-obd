package main

import (
	"fmt"
	cartservice_server_rest_server "github.com/kurtosis-tech/new-obd/src/cartservice/api/http_rest/server"
	"github.com/kurtosis-tech/new-obd/src/cartservice/cartstore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net"
	"os"
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
	echoRouter.Use(middleware.Logger())

	echoRouter.Use(TraceIDMiddleware)

	// CORS configuration
	echoRouter.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: defaultCORSOrigins,
		AllowHeaders: defaultCORSHeaders,
	}))

	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	db, err := cartstore.NewDb(dbHost, dbUsername, dbPassword, dbName, dbPort)
	if err != nil {
		logrus.Fatal(err)
	}

	server := NewServer(db)

	cartservice_server_rest_server.RegisterHandlers(echoRouter, cartservice_server_rest_server.NewStrictHandler(server, nil))

	echoRouter.Start(net.JoinHostPort(restAPIHostIP, fmt.Sprint(restAPIPortAddr)))
}
