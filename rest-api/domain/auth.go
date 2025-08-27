package domain

import (
	"github.com/laurentsbrndn/accounting-app/rest-api/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Logout(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error)
}