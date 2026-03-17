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
		{"Pending to InTransit", StatusPending, StatusInTransit, true},
		{"PickedUp to InTransit", StatusPickedUp, StatusInTransit, false},
		{"InTransit to Delivered", StatusInTransit, StatusDelivered, false},
		{"Pending to Cancelled", StatusPending, StatusCancelled, false},
		{"Delivered to Cancelled", StatusDelivered, StatusCancelled, true},
		{"Cancelled to PickedUp", StatusCancelled, StatusPickedUp, true},
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
