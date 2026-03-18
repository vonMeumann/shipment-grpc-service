package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vonMeumann/shipment-grpc-service/internal/domain"
)

type ShipmentUseCase interface {
	CreateShipment(ctx context.Context, ref, origin, dest, driver, unit string, amount, revenue float64) error
	GetShipment(ctx context.Context, ref string) (*domain.Shipment, error)
	AddStatusEvent(ctx context.Context, ref string, status domain.ShipmentStatus, remarks string) error
	GetShipmentHistory(ctx context.Context, ref string) ([]*domain.StatusEvent, error)
}

type shipmentUseCase struct {
	repo domain.Repository
}

func NewShipmentUseCase(repo domain.Repository) ShipmentUseCase {
	return &shipmentUseCase{repo: repo}
}

func (u *shipmentUseCase) CreateShipment(ctx context.Context, ref, origin, dest, driver, unit string, amount, revenue float64) error {
	slog.InfoContext(ctx, "creating shipment", "ref", ref, "origin", origin, "dest", dest)
	
	shipment := domain.NewShipment(ref, origin, dest, driver, unit, amount, revenue)

	if err := shipment.Validate(); err != nil {
		slog.WarnContext(ctx, "shipment validation failed", "ref", ref, "error", err)
		return err
	}

	if err := u.repo.SaveShipment(ctx, shipment); err != nil {
		slog.ErrorContext(ctx, "failed to save shipment", "ref", ref, "error", err)
		return fmt.Errorf("failed to save shipment: %w", err)
	}

	event := &domain.StatusEvent{
		ID:          uuid.New(),
		ShipmentID:  shipment.ID,
		ShipmentRef: ref,
		Status:      domain.StatusPending,
		Timestamp:   time.Now(),
		Remarks:     "Shipment created",
	}
	if err := u.repo.SaveEvent(ctx, event); err != nil {
		slog.ErrorContext(ctx, "failed to save initial event", "ref", ref, "error", err)
		return fmt.Errorf("failed to save initial event: %w", err)
	}

	slog.InfoContext(ctx, "shipment created successfully", "ref", ref, "id", shipment.ID)
	return nil
}

func (u *shipmentUseCase) GetShipment(ctx context.Context, ref string) (*domain.Shipment, error) {
	slog.DebugContext(ctx, "getting shipment", "ref", ref)
	return u.repo.GetShipment(ctx, ref)
}

func (u *shipmentUseCase) AddStatusEvent(ctx context.Context, ref string, status domain.ShipmentStatus, remarks string) error {
	slog.InfoContext(ctx, "adding status event", "ref", ref, "new_status", string(status))
	
	shipment, err := u.repo.GetShipment(ctx, ref)
	if err != nil {
		slog.ErrorContext(ctx, "shipment not found for status update", "ref", ref, "error", err)
		return fmt.Errorf("failed to get shipment: %w", err)
	}

	oldStatus := string(shipment.CurrentStatus)
	if err := shipment.CanTransitionTo(status); err != nil {
		slog.WarnContext(ctx, "invalid status transition attempt", "ref", ref, "from", oldStatus, "to", string(status), "error", err)
		return err
	}

	shipment.CurrentStatus = status
	if err := u.repo.UpdateShipment(ctx, shipment); err != nil {
		slog.ErrorContext(ctx, "failed to update shipment status", "ref", ref, "error", err)
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	event := &domain.StatusEvent{
		ID:          uuid.New(),
		ShipmentID:  shipment.ID,
		ShipmentRef: ref,
		Status:      status,
		Timestamp:   time.Now(),
		Remarks:     remarks,
	}
	if err := u.repo.SaveEvent(ctx, event); err != nil {
		slog.ErrorContext(ctx, "failed to save status event", "ref", ref, "error", err)
		return fmt.Errorf("failed to save status event: %w", err)
	}

	slog.InfoContext(ctx, "status updated successfully", "ref", ref, "from", oldStatus, "to", string(status))
	return nil
}

func (u *shipmentUseCase) GetShipmentHistory(ctx context.Context, ref string) ([]*domain.StatusEvent, error) {
	slog.DebugContext(ctx, "getting shipment history", "ref", ref)
	return u.repo.GetEvents(ctx, ref)
}
