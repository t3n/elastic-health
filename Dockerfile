FROM alpine:3.9
RUN apk --update --no-cache add ca-certificates
COPY ./build/elastic-health /usr/local/bin/elastic-health
ENTRYPOINT ["/usr/local/bin/elastic-health"]
