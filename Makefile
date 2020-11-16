include .env

## start: Starts the scraper until interrupted.
start:
	go build -o ./bin/main ./cmd/console/ && ./bin/main

## install: Install missing dependencies. Runs `go get` internally.
install:
	go get

## test: Run tests.
test:
	go test *.go