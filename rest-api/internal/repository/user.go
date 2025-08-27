package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/laurentsbrndn/accounting-app/rest-api/domain"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

// DeleteRefreshToken implements domain.UserRepository.
func (u *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	dataset := goqu.Delete("refresh_tokens").Where(goqu.Ex{
		"token": token,
	})

	query, args, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args...)
	return err
}

// IsRefreshTokenValid implements domain.UserRepository.
func (u *userRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, string, error) {
	dataset := goqu.From("refresh_tokens").Select("user_id", "expires_at").Where(goqu.Ex{
		"token": token,
	})

	query, args, err := dataset.ToSQL()
	if err != nil {
		return false, "", err
	}

	row := u.db.QueryRowContext(ctx, query, args...)

	var userId string
	var expiresAt time.Time
	if err := row.Scan(&userId, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	if time.Now().After(expiresAt) {
		return false, "", nil
	}

	return true, userId, nil
}

// StoreRefreshToken implements domain.UserRepository.
func (u *userRepository) StoreRefreshToken(ctx context.Context, Id string, token string, exp time.Time) error {
	dataset := goqu.Insert("refresh_tokens").Rows(
		goqu.Record{
			"token":      token,
			"user_id":    Id,
			"expires_at": exp,
		},
	)

	query, args, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args...)
	return err
}

// Create implements domain.UserRepository.
func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := u.db.Insert("users").Rows(user).Executor().ExecContext(ctx)
	return err
}

// FindByIdentifier implements domain.UserRepository.
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
