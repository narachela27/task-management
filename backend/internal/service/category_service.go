package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/model"
	"task-management-backend/internal/repository"
)

var (
	ErrCategoryNotFound      = errors.New("kategori tidak ditemukan")
	ErrCategoryAlreadyExists = errors.New("nama kategori sudah ada")
)

type CategoryService struct {
	CategoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepo: categoryRepo}
}

func (s *CategoryService) Create(ctx context.Context, userID int64, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	name := strings.TrimSpace(req.Name)

	exists, err := s.CategoryRepo.ExistsByNameAndUserID(ctx, name, userID, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryAlreadyExists
	}

	id, err := s.CategoryRepo.Create(ctx, &model.Category{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	created, err := s.CategoryRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:        created.ID,
		UserID:    created.UserID,
		Name:      created.Name,
		CreatedAt: created.CreatedAt.Format(time.RFC3339),
		UpdatedAt: created.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CategoryService) GetAll(ctx context.Context, userID int64) ([]dto.CategoryResponse, error) {
	categories, err := s.CategoryRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, dto.CategoryResponse{
			ID:        c.ID,
			UserID:    c.UserID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *CategoryService) Update(ctx context.Context, userID, categoryID int64, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	name := strings.TrimSpace(req.Name)

	existing, err := s.CategoryRepo.FindByIDAndUserID(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrCategoryNotFound
	}

	duplicate, err := s.CategoryRepo.ExistsByNameAndUserID(ctx, name, userID, categoryID)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, ErrCategoryAlreadyExists
	}

	if _, err := s.CategoryRepo.Update(ctx, categoryID, userID, name); err != nil {
		return nil, err
	}

	updated, err := s.CategoryRepo.FindByIDAndUserID(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:        updated.ID,
		UserID:    updated.UserID,
		Name:      updated.Name,
		CreatedAt: updated.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CategoryService) Delete(ctx context.Context, userID, categoryID int64) error {
	deleted, err := s.CategoryRepo.Delete(ctx, categoryID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrCategoryNotFound
	}
	return nil
}