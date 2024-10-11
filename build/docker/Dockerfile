# Step 1: Modules caching
FROM golang:1.23 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23 as tester
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    make test

# Step 2: Builder
FROM golang:1.23 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -ldflags="-w -s" -o bin/2dpoint

# Step 3: Final
FROM scratch
COPY --from=builder /app/bin/2dpoint /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./internal/views ./internal/views
COPY ./internal/config/config.yaml ./internal/config/config.yaml
EXPOSE 3000
ENTRYPOINT ["/app"]
