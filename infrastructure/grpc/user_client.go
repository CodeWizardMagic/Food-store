package grpc

import (
	"FoodStore-AdvProg2/proto"
	"context"
	"google.golang.org/grpc"
)

type UserClient struct {
	client proto.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewUserServiceClient(conn)
	return &UserClient{client: client}, conn
}

func (c *UserClient) Register(ctx context.Context, in *proto.RegisterRequest, opts ...grpc.CallOption) (*proto.RegisterResponse, error) {
	return c.client.Register(ctx, in, opts...)
}

func (c *UserClient) Authenticate(ctx context.Context, in *proto.AuthenticateRequest, opts ...grpc.CallOption) (*proto.AuthenticateResponse, error) {
	return c.client.Authenticate(ctx, in, opts...)
}

func (c *UserClient) GetProfile(ctx context.Context, in *proto.GetProfileRequest, opts ...grpc.CallOption) (*proto.GetProfileResponse, error) {
	return c.client.GetProfile(ctx, in, opts...)
}

func (c *UserClient) ValidateToken(ctx context.Context, in *proto.ValidateTokenRequest, opts ...grpc.CallOption) (*proto.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, in, opts...)
}