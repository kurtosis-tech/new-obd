package main

import (
	"fmt"
	"github.com/kurtosis-tech/new-obd/src/libraries/events"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	name = "frontend"
)

type frontendServer struct {
	eventsManager *events.EventsManager
}

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	eventsManager, err := events.CreateEventsManager()
	if err != nil {
		logrus.Errorf("An error occurred initializing events manager! No site events will be received.\nError was: %s", err)
	}

	svc := &frontendServer{
		eventsManager: eventsManager,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/ws", svc.handleWebSocket).Methods(http.MethodGet)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r

	// Start the server
	http.Handle("/", r)
	fmt.Println("Server starting on port 8091...")
	http.ListenAndServe(":8091", handler)
}
