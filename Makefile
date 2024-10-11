.PHONY: build run test
.SILENT: build run
.DEFAULT_GOAL := run

build:
	@go build -o bin/2dpoint

run: build 
	@./bin/2dpoint

air: 
	@go build -o bin/2dpoint.exe

test:
	@go clean -testcache
	@gotestsum --packages="./tests/..." --format testname

docker:
	docker build -f ./build/docker/Dockerfile -t audryus/2dpoint-site .