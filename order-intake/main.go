package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/SanishKumar/supplychain/protos" // import the generated protobuf code
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/orders", createOrderHandler)
	router.Get("/orders/{order_id}", getOrderHandler)

	http.ListenAndServe(":8080", router)
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req pb.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// TODO: Validate req.Items (non-empty, etc.)

	// Call Inventory gRPC service to reserve stock
	conn, err := grpc.DialContext(ctx, "inventory-checker:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "failed to connect to inventory service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	invClient := pb.NewInventoryServiceClient(conn)

	res, err := invClient.ReserveStock(ctx, &pb.ReserveStockRequest{OrderId: req.Items[0].Sku})
	if err != nil {
		http.Error(w, "inventory reservation failed", http.StatusInternalServerError)
		return
	}

	// Respond with OrderCreated event (simplified)
	created := &pb.OrderCreated{
		OrderId: "order-" + generateID(),
		Items:   req.Items,
	}
	resp := &pb.CreateOrderResponse{Order: created}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	orderID := chi.URLParam(r, "order_id")
	// TODO: Fetch order info (stubbed below)
	order := &pb.OrderCreated{
		OrderId: orderID,
		Items:   []*pb.OrderItem{}, // stub
	}
	resp := &pb.GetOrderResponse{Order: order}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
