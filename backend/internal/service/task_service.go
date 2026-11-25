package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/model"
	"task-management-backend/internal/repository"
)

var (
	ErrTaskNotFound        = errors.New("task tidak ditemukan")
	ErrInvalidPriority     = errors.New("priority harus bernilai LOW, MEDIUM, atau HIGH")
	ErrInvalidStatus       = errors.New("status harus bernilai TODO, IN_PROGRESS, atau DONE")
	ErrInvalidDueDate      = errors.New("format due_date harus YYYY-MM-DD")
	ErrCategoryNotOwned    = errors.New("kategori tidak ditemukan atau bukan milik Anda")
)

type TaskService struct {
	TaskRepo     *repository.TaskRepository
	CategoryRepo *repository.CategoryRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, categoryRepo *repository.CategoryRepository) *TaskService {
	return &TaskService{
		TaskRepo:     taskRepo,
		CategoryRepo: categoryRepo,
	}
}

func isValidPriority(p string) bool {
	return p == "LOW" || p == "MEDIUM" || p == "HIGH"
}

func isValidStatus(s string) bool {
	return s == "TODO" || s == "IN_PROGRESS" || s == "DONE"
}

func toTaskResponse(t *model.Task) dto.TaskResponse {
	var dueDateStr *string
	if t.DueDate != nil {
		formatted := t.DueDate.Format("2006-01-02")
		dueDateStr = &formatted
	}

	return dto.TaskResponse{
		ID:           t.ID,
		UserID:       t.UserID,
		CategoryID:   t.CategoryID,
		CategoryName: t.CategoryName,
		Title:        t.Title,
		Description:  t.Description,
		Priority:     t.Priority,
		Status:       t.Status,
		DueDate:      dueDateStr,
		CreatedAt:    t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    t.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *TaskService) Create(ctx context.Context, userID int64, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	title := strings.TrimSpace(req.Title)

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if priority == "" {
		priority = "MEDIUM"
	}
	if !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = "TODO"
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	if req.CategoryID != nil {
		cat, err := s.CategoryRepo.FindByIDAndUserID(ctx, *req.CategoryID, userID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, ErrCategoryNotOwned
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil && strings.TrimSpace(*req.DueDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.DueDate))
		if err != nil {
			return nil, ErrInvalidDueDate
		}
		dueDate = &parsed
	}

	id, err := s.TaskRepo.Create(ctx, &model.Task{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       title,
		Description: req.Description,
		Priority:    priority,
		Status:      status,
		DueDate:     dueDate,
	})
	if err != nil {
		return nil, err
	}

	created, err := s.TaskRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	res := toTaskResponse(created)
	return &res, nil
}

func (s *TaskService) GetAll(ctx context.Context, userID int64, filter dto.TaskFilterQuery) ([]dto.TaskResponse, dto.PaginationMeta, error) {
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

	tasks, totalItems, err := s.TaskRepo.FindAllByUserID(ctx, userID, filter)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	items := make([]dto.TaskResponse, 0, len(tasks))
	for i := range tasks {
		items = append(items, toTaskResponse(&tasks[i]))
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	meta := dto.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return items, meta, nil
}

func (s *TaskService) GetByID(ctx context.Context, userID, taskID int64) (*dto.TaskResponse, error) {
	task, err := s.TaskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	res := toTaskResponse(task)
	return &res, nil
}

func (s *TaskService) Update(ctx context.Context, userID, taskID int64, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	existing, err := s.TaskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTaskNotFound
	}

	title := strings.TrimSpace(req.Title)

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if priority == "" {
		priority = existing.Priority
	}
	if !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = existing.Status
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	if req.CategoryID != nil {
		cat, err := s.CategoryRepo.FindByIDAndUserID(ctx, *req.CategoryID, userID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, ErrCategoryNotOwned
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil && strings.TrimSpace(*req.DueDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.DueDate))
		if err != nil {
			return nil, ErrInvalidDueDate
		}
		dueDate = &parsed
	}

	updatedModel := &model.Task{
		ID:          taskID,
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       title,
		Description: req.Description,
		Priority:    priority,
		Status:      status,
		DueDate:     dueDate,
	}

	if _, err := s.TaskRepo.Update(ctx, updatedModel); err != nil {
		return nil, err
	}

	updated, err := s.TaskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	res := toTaskResponse(updated)
	return &res, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, userID, taskID int64, req dto.UpdateTaskStatusRequest) (*dto.TaskResponse, error) {
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	updatedOk, err := s.TaskRepo.UpdateStatus(ctx, taskID, userID, status)
	if err != nil {
		return nil, err
	}
	if !updatedOk {
		return nil, ErrTaskNotFound
	}

	updated, err := s.TaskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	res := toTaskResponse(updated)
	return &res, nil
}

func (s *TaskService) Delete(ctx context.Context, userID, taskID int64) error {
	deleted, err := s.TaskRepo.Delete(ctx, taskID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrTaskNotFound
	}
	return nil
}

func (s *TaskService) DeleteByAdmin(ctx context.Context, taskID int64) error {
	deleted, err := s.TaskRepo.DeleteByID(ctx, taskID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrTaskNotFound
	}
	return nil
}

func (s *TaskService) GetSummary(ctx context.Context, userID int64) (*dto.TaskSummaryResponse, error) {
	return s.TaskRepo.GetSummaryByUserID(ctx, userID)
}