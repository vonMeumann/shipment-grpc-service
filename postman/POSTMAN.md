# Postman gRPC Documentation

This document explains how to use the provided Postman collection to test the Shipment Tracking gRPC service.

## Setup
1. **link**: https://gold-station-792407.postman.co/workspace/namae~8835573a-dd94-48c9-a65c-8cd4761b7f73/collection/69b9d86ffd7a2abd4befa4d9?action=share&creator=20954863
2. **Import Proto**:
   - In any gRPC request, go to the **Service** (or **Service definition**) tab.
   - Click **Import a .proto file**.
   - Select `api/proto/shipment/v1/shipment.proto`.
   - This allows Postman to understand the gRPC methods and validate your JSON messages.

## Available Requests

### 1. CreateShipment
Creates a new shipment in the `PENDING` state.
- **Method**: `shipment.v1.ShipmentService/CreateShipment`
- **Unique Constraint**: `reference_number` must not already exist.
- **Validation**:
  - `reference_number`: 3-50 chars.
  - `origin`, `destination`, `driver_details`, `unit_details`: 2-100 chars.
  - `shipment_amount`, `driver_revenue`: Non-negative numbers.
- **Negative Case**: (Try creating a shipment with an existing `reference_number` or set `shipment_amount` to `-100.0`)

### 2. GetShipment
Retrieves current details of a shipment by its reference number.
- **Method**: `shipment.v1.ShipmentService/GetShipment`
- **Negative Case**: (Use a `reference_number` that does not exist like `NON-EXISTENT-999`)

### 3. AddShipmentEvent
Updates the status of a shipment and adds a record to its history.
- **Method**: `shipment.v1.ShipmentService/AddShipmentEvent`
- **Allowed Statuses**: `PENDING`, `PICKED_UP`, `IN_TRANSIT`, `DELIVERED`, `CANCELLED`.
- **Validation**:
  - Transition must follow the state machine (e.g., `PENDING -> PICKED_UP`).
  - Cannot update status of `DELIVERED` or `CANCELLED` shipments.
- **Negative Case**: (Try to transition from `PENDING` directly to `DELIVERED` or use a delivered shipment)

### 4. GetShipmentHistory
Returns all status events for a specific shipment, sorted by time.
- **Method**: `shipment.v1.ShipmentService/GetShipmentHistory`
- **Negative Case**: (Use a `reference_number` that does not exist to get an empty list `[]`)

## Troubleshooting
- **Connection Error**: Ensure the server is running on `localhost:50051`.
- **Method Not Found**: Make sure you have imported the `.proto` file as described in the Setup section.
- **Validation Error**: Check the server logs (or Postman response) for specific field validation failures (e.g., "reference number too short").
