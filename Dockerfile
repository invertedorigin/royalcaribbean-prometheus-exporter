FROM golang:1.17-alpine AS builder

WORKDIR /go/src/exporter-go

COPY go.sum go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build .

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/exporter-go/royalcaribbean-prometheus-exporter .
ENTRYPOINT ["./royalcaribbean-prometheus-exporter"]