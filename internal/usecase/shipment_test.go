package usecase

import (
	"testing"

	"github.com/VonMeumann/shipment-service/internal/domain"
	"github.com/VonMeumann/shipment-service/internal/infrastructure/repository"
)

func TestShipmentUseCase_AddStatusEvent(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	uc := NewShipmentUseCase(repo)

	ref := "SHP-001"
	err := uc.CreateShipment(ref, "New York", "Los Angeles", "John Doe", "Truck-01", 1000.0, 800.0)
	if err != nil {
		t.Fatalf("failed to create shipment: %v", err)
	}

	// Valid transition
	err = uc.AddStatusEvent(ref, domain.StatusPickedUp, "Picked up from warehouse")
	if err != nil {
		t.Errorf("expected no error for valid transition, got %v", err)
	}

	// Invalid transition
	err = uc.AddStatusEvent(ref, domain.StatusDelivered, "Delivered early")
	if err == nil {
		t.Error("expected error for invalid transition (PickedUp to Delivered), got nil")
	}

	// Verify history
	history, err := uc.GetShipmentHistory(ref)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("expected 2 history events, got %d", len(history))
	}
}
