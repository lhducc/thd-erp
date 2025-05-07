.PHONY: run build test

run:
	go run ./cmd/api/main.go

build:
	mkdir -p build
	go build -o ./build/app ./cmd/api/main.go

test:
	go test ./...
