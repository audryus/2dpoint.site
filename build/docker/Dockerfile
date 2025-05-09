FROM golang:1.24 AS builder

FROM scratch
COPY ./bin/2dpoint /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./internal/views ./internal/views
COPY ./internal/config/config.yaml ./internal/config/config.yaml
EXPOSE 3000
ENTRYPOINT ["/app"]
