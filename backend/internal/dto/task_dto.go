package dto

type CreateTaskRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description *string `json:"description"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description *string `json:"description"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type TaskFilterQuery struct {
	Status     string
	Priority   string
	CategoryID int64
	Search     string
	Page       int
	Limit      int
}

type TaskResponse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	Priority     string  `json:"priority"`
	Status       string  `json:"status"`
	DueDate      *string `json:"due_date"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type TaskSummaryResponse struct {
	Total      int `json:"total"`
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}