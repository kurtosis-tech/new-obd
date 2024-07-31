module github.com/kurtosis-tech/new-obd/src/frontend

go 1.21.9

toolchain go1.22.4

//replace (
//	github.com/kurtosis-tech/new-obd/src/currencyexternalapi => ../currencyexternalapi
//	github.com/kurtosis-tech/new-obd/src/cartservice => ../cartservice
//)

require (
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/kurtosis-tech/new-obd/src/cartservice v0.0.0-20240730221233-4a5135ef7421
	github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi v0.0.0-20240729190836-e99f63c9b42e
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.36.3
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	go.opentelemetry.io/otel v1.11.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.42.0 // indirect
)
