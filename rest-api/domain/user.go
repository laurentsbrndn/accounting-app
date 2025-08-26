package domain

import "context"

type User struct {
	Id          string  `db:"id"`
	Email       string  `db:"email"`
	DisplayName string  `db:"display_name"`
	Username    string  `db:"username"`
	Password    string  `db:"password"`
	Balance     float64 `db:"balance"`
}

// type UserRepository interface {
// 	FindByEmail(ctx context.Context, email string) (User, error)
// }

type UserRepository interface {
	FindByIdentifier(ctx context.Context, identifier string) (User, error)
	Create(ctx context.Context, user User) error
}
