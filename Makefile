.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	@go build -o bin/2dpoint

run: build 
	@./bin/2dpoint

test:
	@go test -v ./tests/...
