# Shipment Tracking gRPC Microservice

This microservice provides a gRPC API for managing shipments and tracking their status history, built with Go and gRPC.

## 1. Instructions on How to Run the Service

### Option 1: Using Docker (Recommended)
This will build and start the server on port `50051` inside a container.

**Start the service**:
```bash
docker-compose up --build
```

**Stop the service**:
```bash
docker-compose down
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
   *Note: Use `SERVER_PORT=8080 go run ...` to change the default port.*

### Option 3: Postman (GUI)
1. Open Postman and click **Import**.
2. Select the file: `postman/Shipment_Service.postman_collection.json`.
3. The requests are pre-configured with sample JSON data.
4. *Note: You may need to select the `.proto` file in the request settings if Postman asks for it.*

---

## 2. Instructions on How to Run the Tests
The project includes unit tests for domain logic and full integration tests for the gRPC handlers.

### Run all tests:
```bash
make test
# OR
go test -v ./...
```

---

## 3. Architecture Overview
The project follows **Clean/Hexagonal Architecture** principles to ensure the core business logic is isolated from external frameworks and infrastructure.

- **Domain Layer (`internal/domain`)**: Contains the `Shipment` entity and the core business rules (state machine for status transitions). It is independent of all other layers.
- **UseCase Layer (`internal/usecase`)**: Orchestrates business flows (e.g., "Create a shipment and log the initial event"). It communicates with the infrastructure via interfaces.
- **Infrastructure Layer (`internal/infrastructure`)**: Implements the technical details:
    - **gRPC Handler**: Translates Protobuf requests to domain calls.
    - **Repository**: An in-memory, thread-safe data store using `sync.RWMutex`.
    - **Logging**: Structured JSON logging using the Go 1.21+ `slog` package.

---

## 4. Explanation of Design Decisions
- **gRPC over REST**: Chosen for high performance, built-in service contracts via Protobuf, and efficient binary serialization.
- **State Machine for Lifecycle**: Status transitions are not just string updates; they are validated by the `CanTransitionTo` method in the domain model. This centralizes business logic and prevents illegal states.
- **In-Memory Repository**: To satisfy the requirement of not depending on a specific database while remaining fully functional for review. It uses a Mutex to handle concurrent gRPC requests safely.
- **Structured Logging (`slog`)**: JSON logs are used to demonstrate production-readiness, allowing for easy parsing by centralized logging systems (ELK, Datadog).
- **Graceful Shutdown**: The server listens for termination signals (`SIGINT`, `SIGTERM`) to finish active requests before stopping, ensuring data integrity.

---

## 5. Any Assumptions Made
- **Initial Status**: Every new shipment is assumed to start in the `PENDING` state.
- **Reference Number**: Assumed the client provides a unique `reference_number` (string) for each shipment.
- **Cancellation**: Assumed a shipment can be cancelled from any state *except* `DELIVERED`.
- **Currency**: Financial fields (`shipment_amount`, `driver_revenue`) are treated as `float64` for simplicity in this prototype.
- **Audit Log**: Assumed that every status change must create a persistent, immutable event in the history log.

---

## 6. Example gRPC Commands (using grpcurl)

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
