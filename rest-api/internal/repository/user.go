package repository

import (
	"context"
	"database/sql"

	"github.com/laurentsbrndn/accounting-app/rest-api/domain"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

// Create implements domain.UserRepository.
func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := u.db.Insert("users").Rows(user).Executor().ExecContext(ctx)
    return err
}

// FindByEmail implements domain.UserRepository.
// func (u *userRepository) FindByEmail(ctx context.Context, email string) (usr domain.User, err error) {
// 	dataset := u.db.From("users").Where(goqu.C("email").Eq(email))
// 	_, err = dataset.ScanStructContext(ctx, &usr)
// 	return
// }

func (u *userRepository) FindByIdentifier(ctx context.Context, identifier string) (usr domain.User, err error) {
	dataset := u.db.From("users").
		Where(
			goqu.Or(
				goqu.C("email").Eq(identifier),
				goqu.C("username").Eq(identifier),
			),
		)
	_, err = dataset.ScanStructContext(ctx, &usr)
	return
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", con),
	}
}
