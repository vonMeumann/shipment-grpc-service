package handler_test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/VonMeumann/shipment-service/api/proto/shipment/v1"
	"github.com/VonMeumann/shipment-service/internal/infrastructure/grpc/handler"
	"github.com/VonMeumann/shipment-service/internal/infrastructure/repository"
	"github.com/VonMeumann/shipment-service/internal/usecase"
)

func TestShipmentService_Integration(t *testing.T) {
	// 1. Setup gRPC server with bufconn (in-memory connection)
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	repo := repository.NewInMemoryRepository()
	uc := usecase.NewShipmentUseCase(repo)
	h := handler.NewShipmentHandler(uc)

	pb.RegisterShipmentServiceServer(s, h)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	defer s.Stop()

	// 2. Setup client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", 
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewShipmentServiceClient(conn)

	// 3. Execution: Create Shipment
	ref := "INT-TEST-001"
	_, err = client.CreateShipment(ctx, &pb.CreateShipmentRequest{
		ReferenceNumber: ref,
		Origin:          "New York",
		Destination:     "Miami",
		ShipmentAmount:  500.0,
	})
	if err != nil {
		t.Fatalf("failed to create shipment: %v", err)
	}

	// 4. Execution: Update Status
	_, err = client.AddShipmentEvent(ctx, &pb.AddShipmentEventRequest{
		ReferenceNumber: ref,
		Status:          "PICKED_UP",
		Remarks:         "Integration test update",
	})
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	// 5. Verification: Get History
	histResp, err := client.GetShipmentHistory(ctx, &pb.GetShipmentHistoryRequest{ReferenceNumber: ref})
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}

	if len(histResp.Events) != 2 {
		t.Errorf("expected 2 history events, got %d", len(histResp.Events))
	}
}
