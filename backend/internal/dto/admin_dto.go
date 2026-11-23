package dto

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type UpdateUserStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type AdminUserFilterQuery struct {
	Role     string
	IsActive string
	Search   string
	Page     int
	Limit    int
}

type AdminTaskFilterQuery struct {
	UserID   int64
	Status   string
	Priority string
	Search   string
	Page     int
	Limit    int
}

type AdminTaskResponse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	UserName     string  `json:"user_name"`
	UserEmail    string  `json:"user_email"`
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