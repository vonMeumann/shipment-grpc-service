package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/vonMeumann/shipment-grpc-service/api/proto/shipment/v1"
	"github.com/vonMeumann/shipment-grpc-service/internal/infrastructure/grpc/handler"
	"github.com/vonMeumann/shipment-grpc-service/internal/infrastructure/repository"
	"github.com/vonMeumann/shipment-grpc-service/internal/usecase"
)

func main() {
	// 1. Structured Logging Setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Configuration Management
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen", "error", err, "port", port)
		os.Exit(1)
	}

	repo := repository.NewInMemoryRepository()
	uc := usecase.NewShipmentUseCase(repo)
	h := handler.NewShipmentHandler(uc)

	s := grpc.NewServer()
	pb.RegisterShipmentServiceServer(s, h)
	reflection.Register(s)

	// Graceful Shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server listening", "port", port)
		if err := s.Serve(lis); err != nil {
			logger.Error("failed to serve", "error", err)
		}
	}()

	<-stop
	logger.Info("shutting down server...")
	s.GracefulStop()
	logger.Info("server stopped")
}
