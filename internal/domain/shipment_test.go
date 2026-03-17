package domain

import (
	"testing"
)

func TestShipment_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name    string
		current ShipmentStatus
		next    ShipmentStatus
		wantErr bool
	}{
		{"Pending to PickedUp", StatusPending, StatusPickedUp, false},
		{"Pending to InTransit (Skip)", StatusPending, StatusInTransit, true},
		{"Pending to Delivered (Skip)", StatusPending, StatusDelivered, true},
		{"PickedUp to InTransit", StatusPickedUp, StatusInTransit, false},
		{"PickedUp to Pending (Reverse)", StatusPickedUp, StatusPending, true},
		{"InTransit to Delivered", StatusInTransit, StatusDelivered, false},
		{"InTransit to PickedUp (Reverse)", StatusInTransit, StatusPickedUp, true},
		{"Pending to Cancelled", StatusPending, StatusCancelled, false},
		{"PickedUp to Cancelled", StatusPickedUp, StatusCancelled, false},
		{"InTransit to Cancelled", StatusInTransit, StatusCancelled, false},
		{"Delivered to Cancelled (Illegal)", StatusDelivered, StatusCancelled, true},
		{"Cancelled to PickedUp (Illegal)", StatusCancelled, StatusPickedUp, true},
		{"Cancelled to InTransit (Illegal)", StatusCancelled, StatusInTransit, true},
		{"Cancelled to Delivered (Illegal)", StatusCancelled, StatusDelivered, true},
		{"Delivered to Pending (Illegal)", StatusDelivered, StatusPending, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shipment{CurrentStatus: tt.current}
			err := s.CanTransitionTo(tt.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("CanTransitionTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShipment_Validate(t *testing.T) {
	tests := []struct {
		name     string
		shipment *Shipment
		wantErr  bool
		err      error
	}{
		{
			"Valid shipment",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "AA", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU", ShipmentAmount: 100, DriverRevenue: 50},
			false,
			nil,
		},
		{
			"Empty reference",
			&Shipment{ReferenceNumber: "", Origin: "AA", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU"},
			true,
			ErrInvalidReference,
		},
		{
			"Too short reference",
			&Shipment{ReferenceNumber: "SH", Origin: "AA", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU"},
			true,
			ErrInvalidReference,
		},
		{
			"Too long reference",
			&Shipment{ReferenceNumber: "THIS-REFERENCE-NUMBER-IS-DEFINITELY-WAY-TOO-LONG-TO-BE-VALID-IN-OUR-SYSTEM-1234567890", Origin: "AA", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU"},
			true,
			ErrInvalidReference,
		},
		{
			"Empty origin",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU"},
			true,
			ErrInvalidOrigin,
		},
		{
			"Negative amount",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "AA", Destination: "BB", DriverDetails: "DD", UnitDetails: "UU", ShipmentAmount: -1},
			true,
			ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shipment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.err {
				t.Errorf("Validate() error = %v, want %v", err, tt.err)
			}
		})
	}
}
