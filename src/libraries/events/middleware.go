package events

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	ctxKeyUserID = "user-id"
)

var avoidURLsPrefix = []string{"/_", "/static"}

func GetWrapsWithEventsManagerMiddleware(manager *EventsManager) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userId := r.Context().Value(ctxKeyUserID)

			urlToTrack := r.URL

			shouldTrackURL := true

			for _, prefix := range avoidURLsPrefix {
				if strings.HasPrefix(urlToTrack.String(), prefix) {
					shouldTrackURL = false
				}
			}

			if manager != nil && shouldTrackURL {
				eventMsg := fmt.Sprintf("USER: %s VISITED: %s", userId, r.URL)
				err := manager.PublishMessage(eventMsg)
				if err != nil {
					logrus.Errorf("error publishing the page visit message with events manager: %v", err)
				}
			}

			if next != nil && r != nil {
				next.ServeHTTP(w, r)
			}

		})
	}
}
