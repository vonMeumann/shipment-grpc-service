.PHONY: all build test clean generate help

# Variables
BINARY_NAME=shipment-server
PROTO_DIR=api/proto/shipment/v1
PROTO_FILE=$(PROTO_DIR)/shipment.proto

all: test build

## build: Build the server binary
build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

## test: Run all tests
test:
	go test -v ./...

## generate: Regenerate gRPC code from .proto files
generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

## run: Run the server
run:
	go run cmd/server/main.go

## clean: Remove binaries and temporary files
clean:
	rm -rf bin/

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
