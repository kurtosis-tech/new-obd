FROM golang:1.21.9-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/hipstershop

# Build Go binary
COPY ./metrics ./metrics
COPY ./libraries ./libraries
WORKDIR /go/src/hipstershop/metrics
RUN go env -w GOPROXY=https://goproxy.io,direct/
RUN go mod download

RUN go build -o /go/src/hipstershop/metricsbin .

# # Deployment container
FROM alpine:latest

# Install dependencies
RUN apk --update --no-cache add ca-certificates protoc

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/hipstershop/metricsbin /hipstershop/metrics
COPY ./metrics/templates ./templates
COPY ./metrics/static ./static

ENTRYPOINT ["/hipstershop/metrics"]
