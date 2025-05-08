.PHONY: lint run build test

lint:
	golangci-lint run --timeout 5m	

run:
	go run ./cmd/api/main.go

build:
	mkdir -p build
	go build -o ./build/app ./cmd/api/main.go

test:
	go test ./...
