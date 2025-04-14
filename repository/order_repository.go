package repository

import "FoodStore-AdvProg2/domain"

type OrderRepository interface {
	Save(order domain.Order, items []domain.OrderItem) (string, error)
	FindByID(id string) (domain.Order, []domain.OrderItem, error)
	UpdateStatus(id string, status string) error
	FindByUserID(userID string) ([]domain.Order, error)
	FindAll() ([]domain.Order, error)
}
