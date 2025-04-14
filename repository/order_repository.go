package repository

import (
	"FoodStore-AdvProg2/domain"
)

type OrderRepository interface {
	Save(order domain.Order, items []domain.OrderItem) (string, error)
	FindByID(id string) (domain.Order, []domain.OrderItem, error)
	UpdateStatus(orderID, status string) error
	FindByUserID(userID string) ([]domain.Order, error)
	FindAll() ([]domain.Order, error)
	DeleteOrderItemsByProduct(productID string) error
}