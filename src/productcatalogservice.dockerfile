FROM golang:1.21.9-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/hipstershop

# Build Go binary
COPY ./productcatalogservice ./productcatalogservice
WORKDIR /go/src/hipstershop/productcatalogservice
RUN go env -w GOPROXY=https://goproxy.io,direct/
RUN go mod download

RUN go build -o /go/src/hipstershop/productcatalogservicebin .

# Deployment container
FROM scratch

WORKDIR /hipstershop

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

COPY productcatalogservice/data /hipstershop/data/
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/hipstershop/productcatalogservicebin /hipstershop/productcatalogservice

ENTRYPOINT ["/hipstershop/productcatalogservice"]
