module github.com/kurtosis-tech/new-obd/src/frontend

go 1.19

replace github.com/kurtosis-tech/new-obd/src/currencyexternalapi => ../src/currencyexternalapi

require (
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi v0.0.0-20240729190836-e99f63c9b42e
	github.com/sirupsen/logrus v1.7.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.44.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	go.opentelemetry.io/otel v1.18.0 // indirect
	go.opentelemetry.io/otel/metric v1.18.0 // indirect
	go.opentelemetry.io/otel/trace v1.18.0 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.42.0 // indirect
)
