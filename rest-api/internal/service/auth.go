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

// Logout implements domain.AuthService.
func (a authService) Logout(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error) {
	err := a.userRepository.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return dto.LogoutResponse{}, err
	}

	return dto.LogoutResponse{
		Message: "logout successful",
	}, nil
}

// Register implements domain.AuthService.
func (a authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	userByEmail, err := a.userRepository.FindByIdentifier(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return dto.RegisterResponse{}, err
	}
	if userByEmail.Id != "" {
		return dto.RegisterResponse{}, errors.New("email already registered")
	}

	userByUsername, err := a.userRepository.FindByIdentifier(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return dto.RegisterResponse{}, err
	}
	if userByUsername.Id != "" {
		return dto.RegisterResponse{}, errors.New("username already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to hash password")
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

	accessClaim := jwt.MapClaims{
		"id":  newUser.Id,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim)
	accessTokenStr, err := accessToken.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to generate access token")
	}

	refreshToken := uuid.New().String()
	refreshExp := time.Now().Add(time.Duration(a.conf.Jwt.RefreshExp) * time.Minute)
	err = a.userRepository.StoreRefreshToken(ctx, newUser.Id, refreshToken, refreshExp)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	return dto.RegisterResponse{
		Id:           newUser.Id,
		Email:        newUser.Email,
		Username:     newUser.Username,
		DisplayName:  newUser.DisplayName,
		AccessToken:  accessTokenStr,
		RefreshToken: refreshToken,
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

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	accessTokenStr, err := accessToken.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.LoginResponse{}, errors.New("authentication failed")
	}

	refreshToken := uuid.New().String()
	refreshExp := time.Now().Add(time.Duration(a.conf.Jwt.RefreshExp) * time.Minute)
	err = a.userRepository.StoreRefreshToken(ctx, user.Id, refreshToken, refreshExp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuth(cnf *config.Config,
	userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}
