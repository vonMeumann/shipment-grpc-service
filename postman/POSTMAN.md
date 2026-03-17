# Postman gRPC Documentation

This document explains how to use the provided Postman collection to test the Shipment Tracking gRPC service.

## Setup
1. **Import Collection**: Import `postman/Shipment_Service.postman_collection.json` into Postman.
2. **Import Proto**:
   - In any gRPC request, go to the **Service** (or **Service definition**) tab.
   - Click **Import a .proto file**.
   - Select `api/proto/shipment/v1/shipment.proto`.
   - This allows Postman to understand the gRPC methods and validate your JSON messages.

## Available Requests

### 1. CreateShipment
Creates a new shipment in the `PENDING` state.
- **Valid Case**: Use standard JSON with all fields filled.
- **Negative Case (Duplicate Ref)**: Try creating a shipment with an existing `reference_number`. 
  - *Expected Error*: `shipment with this reference number already exists`
- **Negative Case (Invalid Amount)**: Set `shipment_amount` to `-100.0`.
  - *Expected Error*: `shipment amount and driver revenue cannot be negative`

### 2. GetShipment
Retrieves current details of a shipment by its reference number.
- **Valid Case**: Use an existing `reference_number`.
- **Negative Case (Not Found)**: Use a `reference_number` that does not exist (e.g., `NON-EXISTENT-999`).
  - *Expected Error*: `shipment not found`

### 3. AddShipmentEvent
Updates the status of a shipment and adds a record to its history.
- **Valid Case**: Transition from `PENDING` to `PICKED_UP`.
- **Negative Case (Invalid Transition)**: Try to transition from `PENDING` directly to `DELIVERED`.
  - *Expected Error*: `invalid status transition`

### 4. GetShipmentHistory
Returns all status events for a specific shipment, sorted by time.
- **Valid Case**: Use an existing `reference_number`.
- **Negative Case (Empty History)**: Use a `reference_number` that does not exist.
  - *Expected Result*: Returns an empty list of events `[]`.

## Troubleshooting
- **Connection Error**: Ensure the server is running on `localhost:50051`.
- **Method Not Found**: Make sure you have imported the `.proto` file as described in the Setup section.
- **Validation Error**: Check the server logs (or Postman response) for specific field validation failures (e.g., "reference number too short").
