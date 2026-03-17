package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/VonMeumann/shipment-service/api/proto/shipment/v1"
	"github.com/VonMeumann/shipment-service/internal/domain"
	"github.com/VonMeumann/shipment-service/internal/usecase"
)

type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer
	usecase usecase.ShipmentUseCase
}

func NewShipmentHandler(u usecase.ShipmentUseCase) *ShipmentHandler {
	return &ShipmentHandler{usecase: u}
}

func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	err := h.usecase.CreateShipment(
		req.ReferenceNumber,
		req.Origin,
		req.Destination,
		req.DriverDetails,
		req.UnitDetails,
		req.ShipmentAmount,
		req.DriverRevenue,
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateShipmentResponse{ReferenceNumber: req.ReferenceNumber}, nil
}

func (h *ShipmentHandler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	shipment, err := h.usecase.GetShipment(req.ReferenceNumber)
	if err != nil {
		return nil, err
	}

	return &pb.GetShipmentResponse{
		Shipment: &pb.Shipment{
			ReferenceNumber: shipment.ReferenceNumber,
			Origin:          shipment.Origin,
			Destination:     shipment.Destination,
			CurrentStatus:   string(shipment.CurrentStatus),
			DriverDetails:   shipment.DriverDetails,
			UnitDetails:     shipment.UnitDetails,
			ShipmentAmount:  shipment.ShipmentAmount,
			DriverRevenue:   shipment.DriverRevenue,
		},
	}, nil
}

func (h *ShipmentHandler) AddShipmentEvent(ctx context.Context, req *pb.AddShipmentEventRequest) (*pb.AddShipmentEventResponse, error) {
	err := h.usecase.AddStatusEvent(req.ReferenceNumber, domain.ShipmentStatus(req.Status), req.Remarks)
	if err != nil {
		return nil, err
	}
	return &pb.AddShipmentEventResponse{Success: true}, nil
}

func (h *ShipmentHandler) GetShipmentHistory(ctx context.Context, req *pb.GetShipmentHistoryRequest) (*pb.GetShipmentHistoryResponse, error) {
	events, err := h.usecase.GetShipmentHistory(req.ReferenceNumber)
	if err != nil {
		return nil, err
	}

	pbEvents := make([]*pb.StatusEvent, len(events))
	for i, e := range events {
		pbEvents[i] = &pb.StatusEvent{
			ShipmentRef: e.ShipmentRef,
			Status:      string(e.Status),
			Timestamp:   timestamppb.New(e.Timestamp),
			Remarks:     e.Remarks,
		}
	}

	return &pb.GetShipmentHistoryResponse{Events: pbEvents}, nil
}
