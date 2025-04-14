package main

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto"
	"FoodStore-AdvProg2/usecase"
	"context"
	"errors"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	grpcpkg "google.golang.org/grpc"
)

type inventoryServer struct {
	proto.UnimplementedInventoryServiceServer
	uc *usecase.ProductUseCase
}

func NewInventoryServer(uc *usecase.ProductUseCase) *inventoryServer {
	return &inventoryServer{uc: uc}
}

func (s *inventoryServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	product := domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	}
	err := s.uc.Create(product)
	if err != nil {
		return nil, err
	}
	return &proto.CreateProductResponse{Id: product.ID}, nil
}

func (s *inventoryServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	product, err := s.uc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	}, nil
}

func (s *inventoryServer) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	product := domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	}
	err := s.uc.Update(req.Id, product)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateProductResponse{
		Id:    req.Id,
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	}, nil
}

func (s *inventoryServer) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	err := s.uc.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteProductResponse{Success: true}, nil
}

func (s *inventoryServer) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	filter := domain.FilterParams{
		Name:     req.Filter.Name,
		MinPrice: req.Filter.MinPrice,
		MaxPrice: req.Filter.MaxPrice,
	}
	pagination := domain.PaginationParams{
		Page:    int(req.Pagination.Page),
		PerPage: int(req.Pagination.PerPage),
	}

	products, total, err := s.uc.List(filter, pagination)
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*proto.Product, len(products))
	for i, p := range products {
		protoProducts[i] = &proto.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		}
	}

	return &proto.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Page:     int32(pagination.Page),
		PerPage:  int32(pagination.PerPage),
	}, nil
}

func (s *inventoryServer) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error) {
	product, err := s.uc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	newStock := product.Stock
	if req.Decrement {
		newStock -= int(req.Stock)
	} else {
		newStock += int(req.Stock)
	}

	if newStock < 0 {
		return nil, errors.New("insufficient stock")
	}

	updatedProduct := domain.Product{
		Name:  product.Name,
		Price: product.Price,
		Stock: newStock,
	}

	err = s.uc.Update(req.Id, updatedProduct)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateStockResponse{Success: true}, nil
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
	log.Println("Connected to PostgreSQL")

	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	productRepo := postgres.NewProductPostgresRepo()
	uc := usecase.NewProductUseCase(productRepo)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpcpkg.NewServer()
	proto.RegisterInventoryServiceServer(grpcServer, NewInventoryServer(uc))

	log.Println("Starting gRPC Inventory Service on :50053...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}