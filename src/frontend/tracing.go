package main

import (
	"github.com/kurtosis-tech/new-obd/src/frontend/consts"
	"github.com/sirupsen/logrus"
	"net/http"
)

func myTracingContextWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceIdFrom := r.Header.Get(consts.KardinalTraceIdHeaderKey)
		logrus.Infof("[KARDINAL-DEBUG] traceIdFrom: %s", traceIdFrom)

		next.ServeHTTP(w, r)
	})
}
