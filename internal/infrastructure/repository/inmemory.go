package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/vonMeumann/shipment-grpc-service/internal/domain"
)

var (
	ErrShipmentNotFound = errors.New("shipment not found")
)

type InMemoryRepository struct {
	mu        sync.RWMutex
	shipments map[string]*domain.Shipment
	events    map[string][]*domain.StatusEvent
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		shipments: make(map[string]*domain.Shipment),
		events:    make(map[string][]*domain.StatusEvent),
	}
}

func (r *InMemoryRepository) SaveShipment(ctx context.Context, s *domain.Shipment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.shipments[s.ReferenceNumber]; exists {
		return domain.ErrShipmentAlreadyExists
	}

	r.shipments[s.ReferenceNumber] = s
	return nil
}

func (r *InMemoryRepository) UpdateShipment(ctx context.Context, s *domain.Shipment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.shipments[s.ReferenceNumber]; !exists {
		return ErrShipmentNotFound
	}

	r.shipments[s.ReferenceNumber] = s
	return nil
}

func (r *InMemoryRepository) GetShipment(ctx context.Context, ref string) (*domain.Shipment, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shipments[ref]
	if !ok {
		return nil, ErrShipmentNotFound
	}
	return s, nil
}

func (r *InMemoryRepository) SaveEvent(ctx context.Context, e *domain.StatusEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[e.ShipmentRef] = append(r.events[e.ShipmentRef], e)
	return nil
}

func (r *InMemoryRepository) GetEvents(ctx context.Context, ref string) ([]*domain.StatusEvent, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	events, ok := r.events[ref]
	if !ok {
		return []*domain.StatusEvent{}, nil
	}
	return events, nil
}
