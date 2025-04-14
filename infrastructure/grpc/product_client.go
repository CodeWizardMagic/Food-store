package grpc

import (
	"FoodStore-AdvProg2/proto"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	client proto.InventoryServiceClient
}

func NewProductClient(addr string) (*ProductClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewInventoryServiceClient(conn)
	return &ProductClient{client: client}, conn
}

func (c *ProductClient) CreateProduct(ctx context.Context, in *proto.CreateProductRequest, opts ...grpc.CallOption) (*proto.CreateProductResponse, error) {
	return c.client.CreateProduct(ctx, in, opts...)
}

func (c *ProductClient) DeleteProduct(ctx context.Context, in *proto.DeleteProductRequest, opts ...grpc.CallOption) (*proto.DeleteProductResponse, error) {
	return c.client.DeleteProduct(ctx, in, opts...)
}

func (c *ProductClient) ListProducts(ctx context.Context, in *proto.ListProductsRequest, opts ...grpc.CallOption) (*proto.ListProductsResponse, error) {
	return c.client.ListProducts(ctx, in, opts...)
}

func (c *ProductClient) UpdateProduct(ctx context.Context, in *proto.UpdateProductRequest, opts ...grpc.CallOption) (*proto.UpdateProductResponse, error) {
	return c.client.UpdateProduct(ctx, in, opts...)
}

func (c *ProductClient) GetProduct(ctx context.Context, in *proto.GetProductRequest, opts ...grpc.CallOption) (*proto.GetProductResponse, error) {
	return c.client.GetProduct(ctx, in, opts...)
}

func (c *ProductClient) UpdateStock(ctx context.Context, in *proto.UpdateStockRequest, opts ...grpc.CallOption) (*proto.UpdateStockResponse, error) {
	return c.client.UpdateStock(ctx, in, opts...)
}