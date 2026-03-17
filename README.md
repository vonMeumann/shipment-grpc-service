# Shipment Tracking gRPC Microservice

This microservice provides a gRPC API for managing shipments and tracking their status history, built with Go and gRPC.

## Architecture

The project follows **Clean/Hexagonal Architecture** principles to ensure separation of concerns and testability:
- **Domain Layer**: Core business models, status lifecycle, and transition rules.
- **UseCase Layer**: Orchestrates business logic and persists data through an abstract repository.
- **Infrastructure Layer**: Implements gRPC handlers, in-memory repository, and configuration.

## Features & Improvements
- [x] **Clean Architecture**: Domain logic is independent of frameworks and databases.
- [x] **Strict Lifecycle Enforcement**: Prevents invalid status transitions (e.g., from `PENDING` to `DELIVERED` without `PICKED_UP`).
- [x] **Dockerized**: Easy to run with `docker-compose`.
- [x] **Structured Logging**: Uses `slog` for JSON logging, suitable for production.
- [x] **Configuration Management**: Supports environment variables (e.g., `SERVER_PORT`).
- [x] **Comprehensive Testing**: Includes unit tests for domain logic and full integration tests.
- [x] **Graceful Shutdown**: Handles OS signals to stop the server cleanly.

## Requirements
- Go 1.22+
- Docker (optional)
- `protoc` (only if regenerating code)

## Getting Started

### Option 1: Using Docker (Recommended)
```bash
docker-compose up --build
```

### Option 2: Running Locally
1. **Install dependencies**:
   ```bash
   go mod download
   ```
2. **Run the server**:
   ```bash
   make run
   # OR
   go run cmd/server/main.go
   ```
   *Note: Use `SERVER_PORT=8080 go run ...` to change the default port (50051).*

## Running Tests
Run both unit and integration tests:
```bash
make test
# OR
go test -v ./...
```

## Status Lifecycle
The service enforces the following status transitions:
- `PENDING` -> `PICKED_UP`
- `PICKED_UP` -> `IN_TRANSIT`
- `IN_TRANSIT` -> `DELIVERED`
- Any status -> `CANCELLED` (except `DELIVERED`)

## Example gRPC Commands (using grpcurl)

### Create a Shipment
```bash
grpcurl -plaintext -d '{
  "reference_number": "SHP-001",
  "origin": "New York",
  "destination": "Los Angeles",
  "driver_details": "John Doe",
  "unit_details": "Truck-01",
  "shipment_amount": 1200.50,
  "driver_revenue": 950.00
}' localhost:50051 shipment.v1.ShipmentService/CreateShipment
```

### Add a Status Event
```bash
grpcurl -plaintext -d '{
  "reference_number": "SHP-001",
  "status": "PICKED_UP",
  "remarks": "Package collected"
}' localhost:50051 shipment.v1.ShipmentService/AddShipmentEvent
```

### Get Shipment Details
```bash
grpcurl -plaintext -d '{"reference_number": "SHP-001"}' localhost:50051 shipment.v1.ShipmentService/GetShipment
```

### Get Shipment History
```bash
grpcurl -plaintext -d '{"reference_number": "SHP-001"}' localhost:50051 shipment.v1.ShipmentService/GetShipmentHistory
```
