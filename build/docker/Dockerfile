# Step 1: Modules caching
FROM golang:latest as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:latest as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 make build

# Step 3: Final
FROM scratch
COPY --from=builder /app/bin/2dpoint /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./internal/views /internal/views
COPY config.yaml config.yaml
EXPOSE 3000
ENTRYPOINT ["/app"]
