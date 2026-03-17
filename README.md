# Shipment Tracking gRPC Microservice

This microservice provides a gRPC API for managing shipments and tracking their status history.

## Architecture

The project follows **Clean/Hexagonal Architecture** principles:
- **Domain Layer**: Contains business models, status lifecycle, and core logic.
- **UseCase Layer**: Orchestrates business operations and enforces rules.
- **Infrastructure Layer**: Implements gRPC handlers and an in-memory repository.

## Requirements
- Go 1.22+
- `protoc` (to regenerate code if needed)

## Getting Started

### 1. Install dependencies
```bash
go mod download
```

### 2. Run the server
```bash
go run cmd/server/main.go
```
The server listens on `:50051`.

### 3. Running tests
```bash
go test ./...
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
