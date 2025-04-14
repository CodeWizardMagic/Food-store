package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(dbHost string) (*pgxpool.Pool, error) {
    log.Printf("Connecting to database with URL: %s", dbHost)

    config, err := pgxpool.ParseConfig(dbHost)
    if err != nil {
        log.Printf("Error parsing database URL: %v", err)
        return nil, err
    }

    config.MaxConns = 10
    config.MinConns = 1
    config.HealthCheckPeriod = 1 * time.Minute

    db, err := pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        log.Printf("Error connecting to database: %v", err)
        return nil, err
    }

    if err := db.Ping(context.Background()); err != nil {
        log.Printf("Error pinging database: %v", err)
        return nil, err
    }

    log.Println("Database connection established")
    return db, nil
}