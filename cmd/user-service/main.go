package main

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto"
	"FoodStore-AdvProg2/usecase"
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type userServer struct {
	proto.UnimplementedUserServiceServer
	uc *usecase.UserUseCase
}

func NewUserServer(uc *usecase.UserUseCase) *userServer {
	return &userServer{uc: uc}
}

func (s *userServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	userID, err := s.uc.Register(user)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{UserId: userID}, nil
}

func (s *userServer) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	token, userID, err := s.uc.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthenticateResponse{
		UserId: userID,
		Token:  token,
	}, nil
}

func (s *userServer) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
	user, err := s.uc.GetProfile(req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.GetProfileResponse{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *userServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	userID, err := s.uc.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &proto.ValidateTokenResponse{UserId: userID}, nil
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

	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	userRepo := postgres.NewUserPostgresRepo(db)
	uc := usecase.NewUserUseCase(userRepo)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, NewUserServer(uc))

	log.Println("Starting gRPC User Service on :50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}