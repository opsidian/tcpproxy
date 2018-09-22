FROM golang:1.11-alpine3.8 AS build

COPY . /go/src/github.com/google/tcpproxy

RUN cd /go/src/github.com/google/tcpproxy && \
    go build -o /bin/dnssrvrouter cmd/dnssrvrouter/main.go && \
    chmod +x /bin/dnssrvrouter

FROM alpine:3.8

COPY --from=build /bin/dnssrvrouter /bin/dnssrvrouter

ENTRYPOINT ["/bin/dnssrvrouter"]
