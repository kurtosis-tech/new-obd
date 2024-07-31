package main

import (
	"github.com/kurtosis-tech/new-obd/src/cartservice/consts"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func myTracingContextWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceIdFrom := r.Header.Get("X-Kardinal-Trace-Id")
		logrus.Infof("[KARDINAL-DEBUG] traceIdFrom: %s", traceIdFrom)

	})
}

// TraceIDMiddleware logs the trace ID from the request headers
func TraceIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the trace ID from the request header
		traceID := c.Request().Header.Get(consts.KardinalTraceIdHeaderKey)

		// Log the trace ID
		if traceID != "" {
			logrus.Infof("[KARDINAL-DEBUG] traceIdFrom: %s", traceID)
		} else {
			log.Println("[KARDINAL-DEBUG] Trace ID: not provided")
		}

		// Call the next handler
		return next(c)
	}
}
