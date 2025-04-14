package usecase

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/proto"
	"FoodStore-AdvProg2/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	orderRepo     repository.OrderRepository
	productClient proto.InventoryServiceClient
	userClient    proto.UserServiceClient
}

func NewOrderUseCase(orderRepo repository.OrderRepository, productClient proto.InventoryServiceClient, userClient proto.UserServiceClient) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:     orderRepo,
		productClient: productClient,
		userClient:    userClient,
	}
}

func (uc *OrderUseCase) CreateOrder(req domain.OrderRequest) (string, error) {
	_, err := uc.userClient.GetProfile(context.Background(), &proto.GetProfileRequest{UserId: req.UserID})
	if err != nil {
		return "", errors.New("invalid user")
	}

	var totalPrice float64
	items := make([]domain.OrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		resp, err := uc.productClient.GetProduct(context.Background(), &proto.GetProductRequest{Id: itemReq.ProductID})
		if err != nil {
			return "", errors.New("invalid product")
		}
		if resp.Stock < int32(itemReq.Quantity) {
			return "", errors.New("insufficient stock")
		}

		items[i] = domain.OrderItem{
			ID:        uuid.New().String(),
			OrderID:   "", 
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     resp.Price,
		}
		totalPrice += float64(itemReq.Quantity) * resp.Price
	}

	order := domain.Order{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		TotalPrice: totalPrice,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	orderID, err := uc.orderRepo.Save(order, items)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		_, err := uc.productClient.UpdateStock(context.Background(), &proto.UpdateStockRequest{
			Id:        item.ProductID,
			Stock:     int32(item.Quantity),
			Decrement: true,
		})
		if err != nil {
			return "", errors.New("failed to update stock")
		}
	}

	return orderID, nil
}

func (uc *OrderUseCase) GetOrderByID(id string) (domain.Order, error) {
	order, items, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}
	order.Items = items
	return order, nil
}

func (uc *OrderUseCase) UpdateOrderStatus(orderID, status string) error {
	return uc.orderRepo.UpdateStatus(orderID, status)
}

func (uc *OrderUseCase) GetOrdersByUserID(userID string) ([]domain.Order, error) {
	return uc.orderRepo.FindByUserID(userID)
}

func (uc *OrderUseCase) GetAllOrders() ([]domain.Order, error) {
	return uc.orderRepo.FindAll()
}
func (uc *OrderUseCase) DeleteOrderItemsByProduct(productID string) error {
    return uc.orderRepo.DeleteOrderItemsByProduct(productID)
}