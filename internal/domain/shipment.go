package domain

import (
	"errors"
	"time"
)

type ShipmentStatus string

const (
	StatusPending   ShipmentStatus = "PENDING"
	StatusPickedUp  ShipmentStatus = "PICKED_UP"
	StatusInTransit ShipmentStatus = "IN_TRANSIT"
	StatusDelivered ShipmentStatus = "DELIVERED"
	StatusCancelled ShipmentStatus = "CANCELLED"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrShipmentAlreadyCancelled = errors.New("shipment already cancelled")
	ErrShipmentAlreadyDelivered = errors.New("shipment already delivered")
)

type Shipment struct {
	ReferenceNumber string
	Origin          string
	Destination     string
	CurrentStatus   ShipmentStatus
	DriverDetails   string
	UnitDetails     string
	ShipmentAmount  float64
	DriverRevenue   float64
}

type StatusEvent struct {
	ShipmentRef string
	Status      ShipmentStatus
	Timestamp   time.Time
	Remarks     string
}

func NewShipment(ref, origin, dest, driver, unit string, amount, revenue float64) *Shipment {
	return &Shipment{
		ReferenceNumber: ref,
		Origin:          origin,
		Destination:     dest,
		CurrentStatus:   StatusPending,
		DriverDetails:   driver,
		UnitDetails:     unit,
		ShipmentAmount:  amount,
		DriverRevenue:   revenue,
	}
}

func (s *Shipment) CanTransitionTo(next ShipmentStatus) error {
	if s.CurrentStatus == StatusCancelled {
		return ErrShipmentAlreadyCancelled
	}
	if s.CurrentStatus == StatusDelivered {
		return ErrShipmentAlreadyDelivered
	}

	if next == StatusCancelled {
		return nil
	}

	switch s.CurrentStatus {
	case StatusPending:
		if next == StatusPickedUp {
			return nil
		}
	case StatusPickedUp:
		if next == StatusInTransit {
			return nil
		}
	case StatusInTransit:
		if next == StatusDelivered {
			return nil
		}
	}

	return ErrInvalidStatusTransition
}
