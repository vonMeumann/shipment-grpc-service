package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/vonMeumann/shipment-grpc-service/internal/domain"
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
	slog.Info("creating shipment", "ref", ref, "origin", origin, "dest", dest)
	shipment := domain.NewShipment(ref, origin, dest, driver, unit, amount, revenue)
	if err := u.repo.SaveShipment(shipment); err != nil {
		slog.Error("failed to save shipment", "ref", ref, "error", err)
		return fmt.Errorf("failed to save shipment: %w", err)
	}

	event := &domain.StatusEvent{
		ShipmentRef: ref,
		Status:      domain.StatusPending,
		Timestamp:   time.Now(),
		Remarks:     "Shipment created",
	}
	if err := u.repo.SaveEvent(event); err != nil {
		slog.Error("failed to save initial event", "ref", ref, "error", err)
		return fmt.Errorf("failed to save initial event: %w", err)
	}

	slog.Info("shipment created successfully", "ref", ref)
	return nil
}

func (u *shipmentUseCase) GetShipment(ref string) (*domain.Shipment, error) {
	slog.Debug("getting shipment", "ref", ref)
	return u.repo.GetShipment(ref)
}

func (u *shipmentUseCase) AddStatusEvent(ref string, status domain.ShipmentStatus, remarks string) error {
	slog.Info("adding status event", "ref", ref, "new_status", string(status))
	shipment, err := u.repo.GetShipment(ref)
	if err != nil {
		slog.Error("shipment not found for status update", "ref", ref, "error", err)
		return fmt.Errorf("failed to get shipment: %w", err)
	}

	oldStatus := string(shipment.CurrentStatus)
	if err := shipment.CanTransitionTo(status); err != nil {
		slog.Warn("invalid status transition attempt", "ref", ref, "from", oldStatus, "to", string(status), "error", err)
		return err
	}

	shipment.CurrentStatus = status
	if err := u.repo.SaveShipment(shipment); err != nil {
		slog.Error("failed to update shipment status", "ref", ref, "error", err)
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	event := &domain.StatusEvent{
		ShipmentRef: ref,
		Status:      status,
		Timestamp:   time.Now(),
		Remarks:     remarks,
	}
	if err := u.repo.SaveEvent(event); err != nil {
		slog.Error("failed to save status event", "ref", ref, "error", err)
		return fmt.Errorf("failed to save status event: %w", err)
	}

	slog.Info("status updated successfully", "ref", ref, "from", oldStatus, "to", string(status))
	return nil
}

func (u *shipmentUseCase) GetShipmentHistory(ref string) ([]*domain.StatusEvent, error) {
	return u.repo.GetEvents(ref)
}
