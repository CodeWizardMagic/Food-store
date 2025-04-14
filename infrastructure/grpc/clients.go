package grpc

import (
	"FoodStore-AdvProg2/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	InventoryClient proto.InventoryServiceClient
	OrderClient     proto.OrderServiceClient
	UserClient      proto.UserServiceClient
	conns           []*grpc.ClientConn
}

func NewClients(inventoryAddr, orderAddr, userAddr string) (*Clients, error) {
	inventoryConn, err := grpc.Dial(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		inventoryConn.Close()
		return nil, err
	}

	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		inventoryConn.Close()
		orderConn.Close()
		return nil, err
	}

	clients := &Clients{
		InventoryClient: proto.NewInventoryServiceClient(inventoryConn),
		OrderClient:     proto.NewOrderServiceClient(orderConn),
		UserClient:      proto.NewUserServiceClient(userConn),
		conns:           []*grpc.ClientConn{inventoryConn, orderConn, userConn},
	}

	return clients, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		conn.Close()
	}
}