package postgres

import (
    "context"
    "errors"
    "FoodStore-AdvProg2/domain"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
)

type UserPostgresRepo struct {
    db *pgxpool.Pool
}

func NewUserPostgresRepo(db *pgxpool.Pool) *UserPostgresRepo {
    return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) Save(user domain.User) (string, error) {
    ctx := context.Background()
    userID := uuid.New().String()
    createdAt := time.Now()

    _, err := r.db.Exec(ctx, `
        INSERT INTO users (id, username, email, password, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
        userID, user.Username, user.Email, user.Password, createdAt)
    if err != nil {
        return "", err
    }

    return userID, nil
}

func (r *UserPostgresRepo) FindByUsername(username string) (domain.User, error) {
    ctx := context.Background()
    var user domain.User

    err := r.db.QueryRow(ctx, `
        SELECT id, username, email, password, created_at
        FROM users
        WHERE username = $1`, username).
        Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err == pgx.ErrNoRows {
        return domain.User{}, errors.New("user not found")
    }
    if err != nil {
        return domain.User{}, err
    }

    return user, nil
}

func (r *UserPostgresRepo) FindByID(id string) (domain.User, error) {
    ctx := context.Background()
    var user domain.User

    err := r.db.QueryRow(ctx, `
        SELECT id, username, email, password, created_at
        FROM users
        WHERE id = $1`, id).
        Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err == pgx.ErrNoRows {
        return domain.User{}, errors.New("user not found")
    }
    if err != nil {
        return domain.User{}, err
    }

    return user, nil
}

func (r *UserPostgresRepo) SaveToken(token domain.Token) error {
    ctx := context.Background()
    _, err := r.db.Exec(ctx, `
        INSERT INTO tokens (user_id, token, created_at)
        VALUES ($1, $2, $3)`,
        token.UserID, token.Token, token.CreatedAt)
    return err
}

func (r *UserPostgresRepo) FindUserIDByToken(token string) (string, error) {
    ctx := context.Background()
    var userID string

    err := r.db.QueryRow(ctx, `
        SELECT user_id
        FROM tokens
        WHERE token = $1`, token).
        Scan(&userID)
    if err == pgx.ErrNoRows {
        return "", errors.New("token not found")
    }
    if err != nil {
        return "", err
    }

    return userID, nil
}