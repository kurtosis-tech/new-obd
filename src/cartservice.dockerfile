FROM golang:1.21.9-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/hipstershop

# Build Go binary
COPY ./cartservice ./cartservice
WORKDIR /go/src/hipstershop/cartservice
RUN go env -w GOPROXY=https://goproxy.io,direct/
RUN go mod download

RUN go build -o /go/src/hipstershop/cartservicebin .

# # Deployment container
FROM alpine:latest

# Install dependencies
RUN apk --update --no-cache add ca-certificates protoc

# These tools are for debuggin the containres
RUN apk add --no-cache \
    bash \
    curl \
    wget \
    net-tools \
    iputils \
    postgresql-client

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/hipstershop/cartservicebin /hipstershop/cartservice

ENTRYPOINT ["/hipstershop/cartservice"]