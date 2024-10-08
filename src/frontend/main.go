package main

import (
	"fmt"
	"github.com/kurtosis-tech/new-obd/src/libraries/events"
	"github.com/kurtosis-tech/new-obd/src/libraries/tracing"
	"net/http"
	"os"
	"time"

	cartservice_rest_client "github.com/kurtosis-tech/new-obd/src/cartservice/api/http_rest/client"
	"github.com/kurtosis-tech/new-obd/src/frontend/currencyexternalservice"
	productcatalogservice_rest_client "github.com/kurtosis-tech/new-obd/src/productcatalogservice/api/http_rest/client"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	name = "frontend"

	defaultCurrency = "USD"
	cookieMaxAge    = 60 * 60 * 48

	cookiePrefix    = "shop_"
	cookieSessionID = cookiePrefix + "session-id"
	cookieCurrency  = cookiePrefix + "currency"

	userID = "0494c5e0-dde0-48fa-a6d8-f7962f5476bf"
)

type ctxKeySessionID struct{}

type frontendServer struct {
	cartService           *cartservice_rest_client.ClientWithResponses
	productCatalogService *productcatalogservice_rest_client.ClientWithResponses
	currencyService       *currencyexternalservice.CurrencyExternalService
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

	cartServiceHost := os.Getenv("CARTSERVICEHOST")
	productCatalogServiceHost := os.Getenv("PRODUCTCATALOGSERVICEHOST")

	cartServiceServer := fmt.Sprintf("http://%s:8090", cartServiceHost)
	productCatalogServiceServer := fmt.Sprintf("http://%s:8070", productCatalogServiceHost)

	cartServiceClient, err := cartservice_rest_client.NewClientWithResponses(cartServiceServer, cartservice_rest_client.WithHTTPClient(&http.Client{}))
	if err != nil {
		logrus.Fatal("An error occurred creating cart service client!\nError was: %s", err)
	}

	productCatalogServiceClient, err := productcatalogservice_rest_client.NewClientWithResponses(productCatalogServiceServer, productcatalogservice_rest_client.WithHTTPClient(&http.Client{}))
	if err != nil {
		logrus.Fatal("An error occurred creating cart service client!\nError was: %s", err)
	}

	apiKey := os.Getenv("JSDELIVRAPIKEY")

	svc := &frontendServer{
		cartService:           cartServiceClient,
		productCatalogService: productCatalogServiceClient,
		currencyService:       currencyexternalservice.CreateService(apiKey),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/product/{id}", svc.productHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart", svc.addToCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/cart", svc.viewCartHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart/empty", svc.emptyCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/setCurrency", svc.setCurrencyHandler).Methods(http.MethodPost)
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r
	handler = &logHandler{log: log, next: handler} // add logging
	handler = ensureSessionID(handler)             // add session ID
	r.Use(tracing.KardinalTracingContextWrapper)

	// add the events manage middleware
	eventsManager, err := events.CreateEventsManager()
	if err != nil {
		logrus.Errorf("An error occurred initializing events manager! No site events will be tracked due this error.\nError was: %s", err)
	}
	r.Use(events.GetWrapsWithEventsManagerMiddleware(eventsManager))

	// Start the server
	http.Handle("/", r)
	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", handler)
}
