package postgres

import (
    "context"
    "fmt"
    "FoodStore-AdvProg2/domain"
    "github.com/google/uuid"
)

type ProductPostgresRepo struct{}

func NewProductPostgresRepo() *ProductPostgresRepo {
    return &ProductPostgresRepo{}
}

func (r *ProductPostgresRepo) Save(product domain.Product) error {
    query := `INSERT INTO products (id, name, price, stock) VALUES ($1, $2, $3, $4)`
    _, err := DB.Exec(context.Background(), query,
        uuid.New().String(), product.Name, product.Price, product.Stock)
    return err
}

func (r *ProductPostgresRepo) FindByID(id string) (domain.Product, error) {
    query := `SELECT id, name, price, stock FROM products WHERE id = $1`
    row := DB.QueryRow(context.Background(), query, id)

    var p domain.Product
    err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
    return p, err
}

func (r *ProductPostgresRepo) Update(id string, product domain.Product) error {
    query := `UPDATE products SET name=$1, price=$2, stock=$3 WHERE id=$4`
    _, err := DB.Exec(context.Background(), query, product.Name, product.Price, product.Stock, id)
    return err
}

func (r *ProductPostgresRepo) Delete(id string) error {
    query := `DELETE FROM products WHERE id=$1`
    _, err := DB.Exec(context.Background(), query, id)
    return err
}

func (r *ProductPostgresRepo) FindAllWithFilter(filter domain.FilterParams, pagination domain.PaginationParams, offset int) ([]domain.Product, int, error) {
    query := `SELECT id, name, price, stock FROM products WHERE 1=1`
    countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`
    args := []interface{}{}
    argCount := 1

    if filter.MinPrice > 0 {
        query += fmt.Sprintf(" AND price >= $%d", argCount)
        countQuery += fmt.Sprintf(" AND price >= $%d", argCount)
        args = append(args, filter.MinPrice)
        argCount++
    }
    if filter.MaxPrice > 0 {
        query += fmt.Sprintf(" AND price <= $%d", argCount)
        countQuery += fmt.Sprintf(" AND price <= $%d", argCount)
        args = append(args, filter.MaxPrice)
        argCount++
    }

    query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argCount, argCount+1)
    args = append(args, pagination.PerPage, offset)

    var total int
    if len(args) > 2 { 
        err := DB.QueryRow(context.Background(), countQuery, args[:len(args)-2]...).Scan(&total)
        if err != nil {
            return nil, 0, err
        }
    } else { 
        err := DB.QueryRow(context.Background(), countQuery).Scan(&total)
        if err != nil {
            return nil, 0, err
        }
    }

    rows, err := DB.Query(context.Background(), query, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var products []domain.Product
    for rows.Next() {
        var p domain.Product
        err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
        if err != nil {
            return nil, 0, err
        }
        products = append(products, p)
    }

    return products, total, nil
}

func (r *ProductPostgresRepo) FindAll() ([]domain.Product, error) {
    products, _, err := r.FindAllWithFilter(domain.FilterParams{}, domain.PaginationParams{PerPage: 1000}, 0)
    return products, err
}
