.PHONY:

build:
	go build -v ./cmd/apiserver

run: build
	./apiserver

test:
	go test -v -timeout 30s ./...

swag-init:
	swag init -g cmd/apiserver/main.go

swag-fmt:
	swag fmt -g cmd/apiserver/main.go


.DEFAULT_GOAL := run