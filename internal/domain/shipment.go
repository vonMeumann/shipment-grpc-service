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
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
	ErrShipmentAlreadyCancelled = errors.New("shipment already cancelled")
	ErrShipmentAlreadyDelivered = errors.New("shipment already delivered")
	ErrInvalidReference         = errors.New("reference number must be between 3 and 50 characters")
	ErrInvalidOrigin            = errors.New("origin must be between 2 and 100 characters")
	ErrInvalidDestination       = errors.New("destination must be between 2 and 100 characters")
	ErrInvalidDriverDetails     = errors.New("driver details must be between 2 and 100 characters")
	ErrInvalidUnitDetails       = errors.New("unit details must be between 2 and 100 characters")
	ErrInvalidAmount            = errors.New("shipment amount and driver revenue cannot be negative")
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

func (s *Shipment) Validate() error {
	if len(s.ReferenceNumber) < 3 || len(s.ReferenceNumber) > 50 {
		return ErrInvalidReference
	}
	if len(s.Origin) < 2 || len(s.Origin) > 100 {
		return ErrInvalidOrigin
	}
	if len(s.Destination) < 2 || len(s.Destination) > 100 {
		return ErrInvalidDestination
	}
	if len(s.DriverDetails) < 2 || len(s.DriverDetails) > 100 {
		return ErrInvalidDriverDetails
	}
	if len(s.UnitDetails) < 2 || len(s.UnitDetails) > 100 {
		return ErrInvalidUnitDetails
	}
	if s.ShipmentAmount < 0 || s.DriverRevenue < 0 {
		return ErrInvalidAmount
	}
	return nil
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
