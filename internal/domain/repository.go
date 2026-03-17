package domain

type Repository interface {
	SaveShipment(s *Shipment) error
	GetShipment(ref string) (*Shipment, error)
	SaveEvent(e *StatusEvent) error
	GetEvents(ref string) ([]*StatusEvent, error)
}
