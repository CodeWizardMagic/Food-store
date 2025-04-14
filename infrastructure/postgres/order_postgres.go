package postgres

import (
	"FoodStore-AdvProg2/domain"
	"context"
	"github.com/google/uuid"
	"time"
)

type OrderPostgresRepo struct{}

func NewOrderPostgresRepo() *OrderPostgresRepo {
	return &OrderPostgresRepo{}
}

func (r *OrderPostgresRepo) Save(order domain.Order, items []domain.OrderItem) (string, error) {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return "", err
	}
	defer tx.Rollback(context.Background())

	orderID := uuid.New().String()

	orderQuery := `
        INSERT INTO orders (id, user_id, total_price, status) 
        VALUES ($1, $2, $3, $4) 
        RETURNING created_at
    `
	var createdAt time.Time
	err = tx.QueryRow(
		context.Background(),
		orderQuery,
		orderID,
		order.UserID,
		order.TotalPrice,
		domain.OrderStatusPending,
	).Scan(&createdAt)
	if err != nil {
		return "", err
	}

	itemQuery := `
        INSERT INTO order_items (id, order_id, product_id, quantity, price) 
        VALUES ($1, $2, $3, $4, $5)
    `
	for _, item := range items {
		_, err := tx.Exec(
			context.Background(),
			itemQuery,
			uuid.New().String(),
			orderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return "", err
	}

	return orderID, nil
}

func (r *OrderPostgresRepo) FindByID(id string) (domain.Order, []domain.OrderItem, error) {
	orderQuery := `
        SELECT id, user_id, total_price, status, created_at 
        FROM orders 
        WHERE id = $1
    `
	var order domain.Order
	err := DB.QueryRow(context.Background(), orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return domain.Order{}, nil, err
	}

	itemsQuery := `
        SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
               p.name, p.stock
        FROM order_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = $1
    `
	rows, err := DB.Query(context.Background(), itemsQuery, id)
	if err != nil {
		return domain.Order{}, nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var productName string
		var productStock int

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&productName,
			&productStock,
		)
		if err != nil {
			return domain.Order{}, nil, err
		}

		item.Product = &domain.Product{
			ID:    item.ProductID,
			Name:  productName,
			Price: item.Price,
			Stock: productStock,
		}

		items = append(items, item)
	}

	return order, items, nil
}

func (r *OrderPostgresRepo) UpdateStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := DB.Exec(context.Background(), query, status, id)
	return err
}

func (r *OrderPostgresRepo) FindByUserID(userID string) ([]domain.Order, error) {
	query := `
        SELECT id, user_id, total_price, status, created_at 
        FROM orders 
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderPostgresRepo) FindAll() ([]domain.Order, error) {
	query := `
        SELECT id, user_id, total_price, status, created_at 
        FROM orders 
        ORDER BY created_at DESC
    `
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
