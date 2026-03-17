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
