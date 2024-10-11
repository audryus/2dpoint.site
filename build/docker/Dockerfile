FROM scratch
COPY /bin/2dpoint /app
COPY ./internal/views ./internal/views
COPY ./internal/config/config.yaml ./internal/config/config.yaml
EXPOSE 3000
ENTRYPOINT ["/app"]
