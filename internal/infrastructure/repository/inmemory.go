package repository

import (
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

func (r *InMemoryRepository) SaveShipment(s *domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipments[s.ReferenceNumber] = s
	return nil
}

func (r *InMemoryRepository) GetShipment(ref string) (*domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shipments[ref]
	if !ok {
		return nil, ErrShipmentNotFound
	}
	return s, nil
}

func (r *InMemoryRepository) SaveEvent(e *domain.StatusEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[e.ShipmentRef] = append(r.events[e.ShipmentRef], e)
	return nil
}

func (r *InMemoryRepository) GetEvents(ref string) ([]*domain.StatusEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	events, ok := r.events[ref]
	if !ok {
		return []*domain.StatusEvent{}, nil
	}
	return events, nil
}
