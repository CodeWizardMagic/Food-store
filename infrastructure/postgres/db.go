package postgres

import (
    "context"
    "log"
    "github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(dbHost string) {
    var err error
    DB, err = pgxpool.Connect(context.Background(), dbHost)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    log.Println("Connected to PostgreSQL")
}