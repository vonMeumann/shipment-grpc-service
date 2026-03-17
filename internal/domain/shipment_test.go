package domain

import (
	"strings"
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
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16", ShipmentAmount: 100, DriverRevenue: 50},
			false,
			nil,
		},
		{
			"Reference too short",
			&Shipment{ReferenceNumber: "A1", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidReference,
		},
		{
			"Reference too long",
			&Shipment{ReferenceNumber: strings.Repeat("A", 51), Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidReference,
		},
		{
			"Origin too short",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "A", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidOrigin,
		},
		{
			"Origin too long",
			&Shipment{ReferenceNumber: "SHP-001", Origin: strings.Repeat("A", 101), Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidOrigin,
		},
		{
			"Destination too short",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "A", DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidDestination,
		},
		{
			"Destination too long",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: strings.Repeat("A", 101), DriverDetails: "John Doe", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidDestination,
		},
		{
			"Driver details too short",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "J", UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidDriverDetails,
		},
		{
			"Driver details too long",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: strings.Repeat("D", 101), UnitDetails: "Volvo FH16"},
			true,
			ErrInvalidDriverDetails,
		},
		{
			"Unit details too short",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "V"},
			true,
			ErrInvalidUnitDetails,
		},
		{
			"Unit details too long",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: strings.Repeat("U", 101)},
			true,
			ErrInvalidUnitDetails,
		},
		{
			"Negative shipment amount",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16", ShipmentAmount: -100.0, DriverRevenue: 50.0},
			true,
			ErrInvalidAmount,
		},
		{
			"Negative driver revenue",
			&Shipment{ReferenceNumber: "SHP-001", Origin: "Astana", Destination: "Almaty", DriverDetails: "John Doe", UnitDetails: "Volvo FH16", ShipmentAmount: 100.0, DriverRevenue: -50.0},
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
