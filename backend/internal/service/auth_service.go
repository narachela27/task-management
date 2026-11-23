package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/model"
	"task-management-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrAccountInactive    = errors.New("akun dinonaktifkan oleh admin")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)

	existing, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyUsed
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     "USER",
		IsActive: true,
	}

	id, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	created, err := s.UserRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        created.ID,
		Name:      created.Name,
		Email:     created.Email,
		Role:      created.Role,
		IsActive:  created.IsActive,
		CreatedAt: created.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	expireHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if err != nil || expireHours <= 0 {
		expireHours = 24
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: signedToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}