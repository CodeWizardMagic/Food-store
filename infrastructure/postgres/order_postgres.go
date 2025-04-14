package postgres

import (
	"FoodStore-AdvProg2/domain"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderPostgresRepo struct {
	db *pgxpool.Pool
}

func NewOrderPostgresRepo() *OrderPostgresRepo {
	return &OrderPostgresRepo{db: DB}
}

func (r *OrderPostgresRepo) Save(order domain.Order, items []domain.OrderItem) (string, error) {
	ctx := context.Background()
	orderID := uuid.New().String()
	createdAt := time.Now()

	var userExists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, order.UserID).Scan(&userExists)
	if err != nil {
		return "", err
	}
	if !userExists {
		return "", errors.New("user does not exist")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, user_id, total_price, status, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		orderID, order.UserID, order.TotalPrice, order.Status, createdAt)
	if err != nil {
		return "", err
	}

	for i := range items {
		var productExists bool
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`, items[i].ProductID).Scan(&productExists)
		if err != nil {
			return "", err
		}
		if !productExists {
			return "", errors.New("product does not exist")
		}

		itemID := uuid.New().String()
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (id, order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4, $5)`,
			itemID, orderID, items[i].ProductID, items[i].Quantity, items[i].Price)
		if err != nil {
			return "", err
		}
		items[i].ID = itemID
		items[i].OrderID = orderID
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return orderID, nil
}

func (r *OrderPostgresRepo) FindByID(id string) (domain.Order, []domain.OrderItem, error) {
	ctx := context.Background()
	var order domain.Order

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, total_price, status, created_at
		FROM orders
		WHERE id = $1`, id).
		Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt)
	if err == pgx.ErrNoRows {
		return domain.Order{}, nil, errors.New("order not found")
	}
	if err != nil {
		return domain.Order{}, nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items
		WHERE order_id = $1`, id)
	if err != nil {
		return domain.Order{}, nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return domain.Order{}, nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return domain.Order{}, nil, err
	}

	return order, items, nil
}

func (r *OrderPostgresRepo) UpdateStatus(id string, status string) error {
	ctx := context.Background()
	result, err := r.db.Exec(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2`, status, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (r *OrderPostgresRepo) FindByUserID(userID string) ([]domain.Order, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, total_price, status, created_at
		FROM orders
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderPostgresRepo) FindAll() ([]domain.Order, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, total_price, status, created_at
		FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
func (r *OrderPostgresRepo) DeleteOrderItemsByProduct(productID string) error {
    ctx := context.Background()
    _, err := r.db.Exec(ctx, `
        DELETE FROM order_items
        WHERE product_id = $1`, productID)
    return err
}