package usecase

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/repository"
	"errors"
)

type OrderUseCase struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
	}
}

func (uc *OrderUseCase) CreateOrder(orderReq domain.OrderRequest) (string, error) {
	if len(orderReq.Items) == 0 {
		return "", errors.New("order have to have at least one item")
	}

	// Подготовка заказа
	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range orderReq.Items {
		if item.Quantity <= 0 {
			return "", errors.New("order item quantity must be greater than zero")
		}

		product, err := uc.ProductRepo.FindByID(item.ProductID)
		if err != nil {
			return "", errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return "", errors.New("not enough stock for product")
		}

		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)

		totalPrice += product.Price * float64(item.Quantity)
	}

	order := domain.Order{
		UserID:     orderReq.UserID,
		TotalPrice: totalPrice,
		Status:     domain.OrderStatusPending,
	}

	orderID, err := uc.OrderRepo.Save(order, orderItems)
	if err != nil {
		return "", err
	}

	for _, item := range orderReq.Items {
		product, _ := uc.ProductRepo.FindByID(item.ProductID)
		product.Stock -= item.Quantity
		uc.ProductRepo.Update(product.ID, product)
	}

	return orderID, nil
}

func (uc *OrderUseCase) GetOrderByID(id string) (domain.Order, error) {
	order, items, err := uc.OrderRepo.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}

	order.Items = items
	return order, nil
}

func (uc *OrderUseCase) UpdateOrderStatus(id string, status string) error {
	if status != domain.OrderStatusPending &&
		status != domain.OrderStatusCompleted &&
		status != domain.OrderStatusCancelled {
		return errors.New("неверный статус заказа")
	}

	order, _, err := uc.OrderRepo.FindByID(id)
	if err != nil {
		return err
	}

	if status == domain.OrderStatusCancelled && order.Status != domain.OrderStatusCancelled {
		_, items, _ := uc.OrderRepo.FindByID(id)
		for _, item := range items {
			product, _ := uc.ProductRepo.FindByID(item.ProductID)
			product.Stock += item.Quantity
			uc.ProductRepo.Update(product.ID, product)
		}
	}

	return uc.OrderRepo.UpdateStatus(id, status)
}

func (uc *OrderUseCase) GetOrdersByUserID(userID string) ([]domain.Order, error) {
	return uc.OrderRepo.FindByUserID(userID)
}

func (uc *OrderUseCase) GetAllOrders() ([]domain.Order, error) {
	return uc.OrderRepo.FindAll()
}
