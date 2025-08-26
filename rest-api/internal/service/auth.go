package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/laurentsbrndn/accounting-app/rest-api/domain"
	"github.com/laurentsbrndn/accounting-app/rest-api/dto"
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

// Register implements domain.AuthService.
func (a authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	userByEmail, err := a.userRepository.FindByIdentifier(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return dto.RegisterResponse{}, err
	}
	if userByEmail.Id != "" {
		return dto.RegisterResponse{}, errors.New("Email already registered")
	}

	userByUsername, err := a.userRepository.FindByIdentifier(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return dto.RegisterResponse{}, err
	}
	if userByUsername.Id != "" {
		return dto.RegisterResponse{}, errors.New("Username already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("Failed to hash password")
	}

	newUser := domain.User{
		Id:          uuid.New().String(),
		Email:       req.Email,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Password:    string(hashedPassword),
		Balance:     0,
	}

	err = a.userRepository.Create(ctx, newUser)
    if err != nil {
        return dto.RegisterResponse{}, err
    }

	claim := jwt.MapClaims{
        "id":  newUser.Id,
        "exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
    if err != nil {
        return dto.RegisterResponse{}, errors.New("failed to generate token")
    }

	return dto.RegisterResponse{
        Id:          newUser.Id,
        Email:       newUser.Email,
        Username:    newUser.Username,
        DisplayName: newUser.DisplayName,
        Token:       tokenStr,
    }, nil
}

// Login implements domain.AuthService.
func (a authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := a.userRepository.FindByIdentifier(ctx, req.Identifier)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if user.Id == "" {
		return dto.LoginResponse{}, errors.New("authentication failed")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, errors.New("authentication failed")
	}

	claim := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.LoginResponse{}, errors.New("authentication failed")
	}

	return dto.LoginResponse{
		Token: tokenStr,
	}, nil
}

func NewAuth(cnf *config.Config,
	userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}
