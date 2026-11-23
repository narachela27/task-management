package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/repository"
)

var (
	ErrInvalidRole        = errors.New("role harus bernilai USER atau ADMIN")
	ErrCannotModifySelf   = errors.New("admin tidak dapat mengubah role atau menonaktifkan akun sendiri")
)

type AdminService struct {
	UserRepo *repository.UserRepository
	TaskRepo *repository.TaskRepository
}

func NewAdminService(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) *AdminService {
	return &AdminService{
		UserRepo: userRepo,
		TaskRepo: taskRepo,
	}
}

func (s *AdminService) GetAllUsers(ctx context.Context, filter dto.AdminUserFilterQuery) ([]dto.UserResponse, dto.PaginationMeta, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	filter.Role = strings.ToUpper(strings.TrimSpace(filter.Role))
	filter.IsActive = strings.ToLower(strings.TrimSpace(filter.IsActive))
	filter.Search = strings.TrimSpace(filter.Search)

	users, total, err := s.UserRepo.FindAllUsers(ctx, filter)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		items = append(items, dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(filter.Limit)))
	}

	return items, dto.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetUserByID(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	u, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AdminService) UpdateUserRole(ctx context.Context, currentAdminID, targetUserID int64, req dto.UpdateUserRoleRequest) (*dto.UserResponse, error) {
	if currentAdminID == targetUserID {
		return nil, ErrCannotModifySelf
	}

	role := strings.ToUpper(strings.TrimSpace(req.Role))
	if role != "USER" && role != "ADMIN" {
		return nil, ErrInvalidRole
	}

	updatedOk, err := s.UserRepo.UpdateRole(ctx, targetUserID, role)
	if err != nil {
		return nil, err
	}
	if !updatedOk {
		return nil, ErrUserNotFound
	}

	return s.GetUserByID(ctx, targetUserID)
}

func (s *AdminService) UpdateUserStatus(ctx context.Context, currentAdminID, targetUserID int64, req dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	if currentAdminID == targetUserID {
		return nil, ErrCannotModifySelf
	}

	updatedOk, err := s.UserRepo.UpdateStatus(ctx, targetUserID, *req.IsActive)
	if err != nil {
		return nil, err
	}
	if !updatedOk {
		return nil, ErrUserNotFound
	}

	return s.GetUserByID(ctx, targetUserID)
}

func (s *AdminService) GetAllTasks(ctx context.Context, filter dto.AdminTaskFilterQuery) ([]dto.AdminTaskResponse, dto.PaginationMeta, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	filter.Status = strings.ToUpper(strings.TrimSpace(filter.Status))
	filter.Priority = strings.ToUpper(strings.TrimSpace(filter.Priority))
	filter.Search = strings.TrimSpace(filter.Search)

	items, total, err := s.TaskRepo.FindAllForAdmin(ctx, filter)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}
	if items == nil {
		items = make([]dto.AdminTaskResponse, 0)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(filter.Limit)))
	}

	return items, dto.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}