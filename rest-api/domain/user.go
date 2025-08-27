package domain

import (
	"context"
	"time"
)

type User struct {
	Id          string  `db:"id"`
	Email       string  `db:"email"`
	DisplayName string  `db:"display_name"`
	Username    string  `db:"username"`
	Password    string  `db:"password"`
	Balance     float64 `db:"balance"`
}

type UserRepository interface {
	FindByIdentifier(ctx context.Context, identifier string) (User, error)
	Create(ctx context.Context, user User) error
	StoreRefreshToken(ctx context.Context, Id, token string, exp time.Time) error
	DeleteRefreshToken(ctx context.Context, token string) error
	IsRefreshTokenValid(ctx context.Context, token string) (bool, string, error)
}
