package postgres

import (
	"context"
	"log"
)

func InitTables() error {
	createProductsTable := `
    CREATE TABLE IF NOT EXISTS products (
        id UUID PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL
    );`

	createOrdersTable := `
    CREATE TABLE IF NOT EXISTS orders (
        id UUID PRIMARY KEY,
        user_id VARCHAR(255) NOT NULL,
        total_price DECIMAL(10, 2) NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	createOrderItemsTable := `
    CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL
    );`

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	createTokensTable := `
    CREATE TABLE IF NOT EXISTS tokens (
        user_id UUID NOT NULL REFERENCES users(id),
        token VARCHAR(255) PRIMARY KEY,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	tables := []string{
		createProductsTable,
		createOrdersTable,
		createOrderItemsTable,
		createUsersTable,
		createTokensTable,
	}

	for _, table := range tables {
		_, err := DB.Exec(context.Background(), table)
		if err != nil {
			log.Printf("Error creating table: %v", err)
			return err
		}
	}

	log.Println("Tables created successfully")
	return nil
}
