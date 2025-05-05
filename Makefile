.DEFAULT_GOAL := build

.PHONY: fmt lint vet build



fmt:
	go fmt ./...

lint: fmt
	staticcheck ./...

vet: fmt

build: vet
	go mod tidy
	go build -ldflags="-s -w" -o build/bot cmd/bot/main.go
	go build -ldflags="-s -w" -o build/deploy cmd/deploy/main.go
