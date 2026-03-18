# Postman gRPC Testing Guide

This document explains how to test the Shipment Tracking gRPC service using the shared Postman collection.

## Setup
1. **Link**: Open the [Shipment Service Collection](https://www.postman.com/smagulalkey/workspace/testtask/collection/8a3d7b42-1e5f-4d2b-a01c-6d9b1c2e3f4g).
2. **Server Address**: Ensure the address is `localhost:50051`.
3. **Import Proto Definition**:
   - In any request, click on the **Service** dropdown or tab.
   - Click **Import a .proto file**.
   - Select the project's file: `api/proto/shipment/v1/shipment.proto`.
   - Choose **Import as API**.
   - This enables autocomplete and validation for gRPC methods.

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
