FROM golang:1.21.9-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/hipstershop

# Install dependencies
RUN apk --update --no-cache add ca-certificates make

# Build Go binary
COPY ./frontend ./frontend
COPY ./cartservice ./cartservice
COPY ./productcatalogservice ./productcatalogservice
COPY ./libraries ./libraries
WORKDIR /go/src/hipstershop/frontend

RUN go env -w GOPROXY=https://goproxy.io,direct/
RUN go mod download

# build binary
RUN go build -o /go/src/hipstershop/fe .

# Deployment container
FROM scratch

COPY --from=builder /go/src/hipstershop/fe /hipstershop/frontend
COPY ./frontend/templates ./templates
COPY ./frontend/static ./static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/hipstershop/frontend"]
