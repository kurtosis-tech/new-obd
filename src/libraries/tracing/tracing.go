package tracing

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func KardinalTracingContextWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := r.Header.Get(KardinalTraceIdHeaderKey)
		logrus.Infof("[KARDINAL-DEBUG] Trace ID: %s", traceID)

		if next != nil && r != nil {
			next.ServeHTTP(w, r)
		}

	})
}
