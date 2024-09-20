module github.com/kurtosis-tech/new-obd/src/metrics

go 1.21.9

toolchain go1.22.4

replace (
	github.com/kurtosis-tech/new-obd/src/libraries/events => ./../libraries/events
)

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/kurtosis-tech/new-obd/src/libraries/events v0.0.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
