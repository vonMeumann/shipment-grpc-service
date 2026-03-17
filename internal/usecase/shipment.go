package usecase

import (
	"fmt"
	"time"

	"github.com/VonMeumann/shipment-service/internal/domain"
)

type ShipmentUseCase interface {
	CreateShipment(ref, origin, dest, driver, unit string, amount, revenue float64) error
	GetShipment(ref string) (*domain.Shipment, error)
	AddStatusEvent(ref string, status domain.ShipmentStatus, remarks string) error
	GetShipmentHistory(ref string) ([]*domain.StatusEvent, error)
}

type shipmentUseCase struct {
	repo domain.Repository
}

func NewShipmentUseCase(repo domain.Repository) ShipmentUseCase {
	return &shipmentUseCase{repo: repo}
}

func (u *shipmentUseCase) CreateShipment(ref, origin, dest, driver, unit string, amount, revenue float64) error {
	shipment := domain.NewShipment(ref, origin, dest, driver, unit, amount, revenue)
	if err := u.repo.SaveShipment(shipment); err != nil {
		return fmt.Errorf("failed to save shipment: %w", err)
	}

	event := &domain.StatusEvent{
		ShipmentRef: ref,
		Status:      domain.StatusPending,
		Timestamp:   time.Now(),
		Remarks:     "Shipment created",
	}
	if err := u.repo.SaveEvent(event); err != nil {
		return fmt.Errorf("failed to save initial event: %w", err)
	}

	return nil
}

func (u *shipmentUseCase) GetShipment(ref string) (*domain.Shipment, error) {
	return u.repo.GetShipment(ref)
}

func (u *shipmentUseCase) AddStatusEvent(ref string, status domain.ShipmentStatus, remarks string) error {
	shipment, err := u.repo.GetShipment(ref)
	if err != nil {
		return fmt.Errorf("failed to get shipment: %w", err)
	}

	if err := shipment.CanTransitionTo(status); err != nil {
		return err
	}

	shipment.CurrentStatus = status
	if err := u.repo.SaveShipment(shipment); err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	event := &domain.StatusEvent{
		ShipmentRef: ref,
		Status:      status,
		Timestamp:   time.Now(),
		Remarks:     remarks,
	}
	if err := u.repo.SaveEvent(event); err != nil {
		return fmt.Errorf("failed to save status event: %w", err)
	}

	return nil
}

func (u *shipmentUseCase) GetShipmentHistory(ref string) ([]*domain.StatusEvent, error) {
	return u.repo.GetEvents(ref)
}
