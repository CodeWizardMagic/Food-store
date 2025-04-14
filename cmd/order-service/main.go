package main

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/grpc"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto"
	"FoodStore-AdvProg2/usecase"
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	grpcpkg "google.golang.org/grpc"
)

type orderServer struct {
	proto.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderServer(uc *usecase.OrderUseCase) *orderServer {
	return &orderServer{uc: uc}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	items := make([]domain.OrderItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItemRequest{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		}
	}

	orderReq := domain.OrderRequest{
		UserID: req.UserId,
		Items:  items,
	}

	orderID, err := s.uc.CreateOrder(orderReq)
	if err != nil {
		return nil, err
	}

	return &proto.CreateOrderResponse{OrderId: orderID}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	order, err := s.uc.GetOrderByID(req.OrderId)
	if err != nil {
		return nil, err
	}

	items := make([]*proto.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = &proto.OrderItem{
			Id:        item.ID,
			OrderId:   item.OrderID,
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &proto.OrderResponse{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt.Unix(),
		Items:      items,
	}, nil
}

func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.UpdateOrderStatusResponse, error) {
	err := s.uc.UpdateOrderStatus(req.OrderId, req.Status)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateOrderStatusResponse{Status: "updated"}, nil
}

func (s *orderServer) GetUserOrders(ctx context.Context, req *proto.GetUserOrdersRequest) (*proto.GetUserOrdersResponse, error) {
	var orders []domain.Order
	var err error

	if req.UserId == "" {
		orders, err = s.uc.GetAllOrders()
	} else {
		orders, err = s.uc.GetOrdersByUserID(req.UserId)
	}
	if err != nil {
		return nil, err
	}

	responseOrders := make([]*proto.OrderResponse, len(orders))
	for i, order := range orders {
		items := make([]*proto.OrderItem, len(order.Items))
		for j, item := range order.Items {
			items[j] = &proto.OrderItem{
				Id:        item.ID,
				OrderId:   item.OrderID,
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
				Price:     item.Price,
			}
		}
		responseOrders[i] = &proto.OrderResponse{
			Id:         order.ID,
			UserId:     order.UserID,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.Unix(),
			Items:      items,
		}
	}

	return &proto.GetUserOrdersResponse{Orders: responseOrders}, nil
}
func (s *orderServer) DeleteOrderItemsByProduct(ctx context.Context, req *proto.DeleteOrderItemsByProductRequest) (*proto.DeleteOrderItemsByProductResponse, error) {
    err := s.uc.DeleteOrderItemsByProduct(req.ProductId)
    if err != nil {
        return nil, err
    }
    return &proto.DeleteOrderItemsByProductResponse{Success: true}, nil
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	dbHost := os.Getenv("DB")
	if dbHost == "" {
		log.Fatal("DB environment variable not set")
	}

	db, err := postgres.InitDB(dbHost)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	postgres.DB = db 
	log.Println("Connected to PostgreSQL via pgxpool")

	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	orderRepo := postgres.NewOrderPostgresRepo()
	userClient, userConn := grpc.NewUserClient("localhost:50052")
	defer userConn.Close()
	productClient, productConn := grpc.NewProductClient("localhost:50053")
	defer productConn.Close()
	uc := usecase.NewOrderUseCase(orderRepo, productClient, userClient)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpcpkg.NewServer()
	proto.RegisterOrderServiceServer(grpcServer, NewOrderServer(uc))

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}