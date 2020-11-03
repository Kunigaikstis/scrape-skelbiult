include .env

## start: Starts the scraper until interrupted.
start:
	go run main.go repository.go

## install: Install missing dependencies. Runs `go get` internally.
install:
	go get

## test: Run tests.
test:
	go test *.go