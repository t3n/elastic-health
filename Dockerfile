FROM golang:1.13-alpine AS build

RUN apk add --update --no-cache ca-certificates git alpine-sdk

COPY . /go/src/github.com/t3n/elastic-health/

RUN cd /go/src/github.com/t3n/elastic-health/ \
&&  go build -a -o build/elastic-health main.go

FROM alpine:3.11.3
COPY --from=build /go/src/github.com/t3n/elastic-health/build/elastic-health /usr/local/bin/elastic-health
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/usr/local/bin/elastic-health"]
