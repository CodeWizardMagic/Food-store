package domain

import (
	"time"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   *Product `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	Price     float64  `json:"price"`
}

type OrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderStatusUpdateRequest struct {
	Status string `json:"status"`
}
