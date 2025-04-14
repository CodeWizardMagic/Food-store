package domain

import "time"

type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
}

type Token struct {
    UserID    string    `json:"user_id"`
    Token     string    `json:"token"`
    CreatedAt time.Time `json:"created_at"`
}