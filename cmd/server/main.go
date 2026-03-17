package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/VonMeumann/shipment-service/api/proto/shipment/v1"
	"github.com/VonMeumann/shipment-service/internal/infrastructure/grpc/handler"
	"github.com/VonMeumann/shipment-service/internal/infrastructure/repository"
	"github.com/VonMeumann/shipment-service/internal/usecase"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewInMemoryRepository()
	uc := usecase.NewShipmentUseCase(repo)
	h := handler.NewShipmentHandler(uc)

	s := grpc.NewServer()
	pb.RegisterShipmentServiceServer(s, h)
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
