.PHONY:

build:
	go build -v ./cmd/apiserver

run: build
	./apiserver

test:
	go test -v -timeout 30s ./...

.DEFAULT_GOAL := run