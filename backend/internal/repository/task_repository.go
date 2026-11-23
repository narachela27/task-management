package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/model"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) (int64, error) {
	query := `INSERT INTO tasks (user_id, category_id, title, description, priority, status, due_date, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	res, err := r.DB.ExecContext(ctx, query,
		task.UserID, task.CategoryID, task.Title, task.Description, task.Priority, task.Status, task.DueDate,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TaskRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*model.Task, error) {
	query := `SELECT t.id, t.user_id, t.category_id, c.name, t.title, t.description, t.priority, t.status, t.due_date, t.created_at, t.updated_at
	          FROM tasks t
	          LEFT JOIN categories c ON t.category_id = c.id
	          WHERE t.id = ? AND t.user_id = ? LIMIT 1`

	var t model.Task
	var catID sql.NullInt64
	var catName, desc sql.NullString
	var dueDate sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &catID, &catName, &t.Title, &desc, &t.Priority, &t.Status, &dueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if catID.Valid {
		t.CategoryID = &catID.Int64
	}
	if catName.Valid {
		t.CategoryName = &catName.String
	}
	if desc.Valid {
		t.Description = &desc.String
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}

	return &t, nil
}

func (r *TaskRepository) FindAllByUserID(ctx context.Context, userID int64, filter dto.TaskFilterQuery) ([]model.Task, int, error) {
	conditions := []string{"t.user_id = ?"}
	args := []interface{}{userID}

	if filter.Status != "" {
		conditions = append(conditions, "t.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Priority != "" {
		conditions = append(conditions, "t.priority = ?")
		args = append(args, filter.Priority)
	}
	if filter.CategoryID > 0 {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(t.title LIKE ? OR t.description LIKE ?)")
		keyword := "%" + filter.Search + "%"
		args = append(args, keyword, keyword)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	var totalItems int
	countQuery := "SELECT COUNT(1) FROM tasks t " + whereClause
	if err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	query := `SELECT t.id, t.user_id, t.category_id, c.name, t.title, t.description, t.priority, t.status, t.due_date, t.created_at, t.updated_at
	          FROM tasks t
	          LEFT JOIN categories c ON t.category_id = c.id ` + whereClause + `
	          ORDER BY t.created_at DESC LIMIT ? OFFSET ?`

	queryArgs := append(args, filter.Limit, offset)
	rows, err := r.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var catID sql.NullInt64
		var catName, desc sql.NullString
		var dueDate sql.NullTime

		if err := rows.Scan(
			&t.ID, &t.UserID, &catID, &catName, &t.Title, &desc, &t.Priority, &t.Status, &dueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if catID.Valid {
			t.CategoryID = &catID.Int64
		}
		if catName.Valid {
			t.CategoryName = &catName.String
		}
		if desc.Valid {
			t.Description = &desc.String
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}

		tasks = append(tasks, t)
	}

	return tasks, totalItems, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) (bool, error) {
	query := `UPDATE tasks
	          SET category_id = ?, title = ?, description = ?, priority = ?, status = ?, due_date = ?, updated_at = NOW()
	          WHERE id = ? AND user_id = ?`

	res, err := r.DB.ExecContext(ctx, query,
		task.CategoryID, task.Title, task.Description, task.Priority, task.Status, task.DueDate, task.ID, task.UserID,
	)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id, userID int64, status string) (bool, error) {
	query := `UPDATE tasks SET status = ?, updated_at = NOW() WHERE id = ? AND user_id = ?`
	res, err := r.DB.ExecContext(ctx, query, status, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (r *TaskRepository) Delete(ctx context.Context, id, userID int64) (bool, error) {
	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`
	res, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (r *TaskRepository) GetSummaryByUserID(ctx context.Context, userID int64) (*dto.TaskSummaryResponse, error) {
	query := `SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'TODO' THEN 1 ELSE 0 END), 0) AS todo,
		COALESCE(SUM(CASE WHEN status = 'IN_PROGRESS' THEN 1 ELSE 0 END), 0) AS in_progress,
		COALESCE(SUM(CASE WHEN status = 'DONE' THEN 1 ELSE 0 END), 0) AS done
	FROM tasks WHERE user_id = ?`

	var summary dto.TaskSummaryResponse
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&summary.Total, &summary.Todo, &summary.InProgress, &summary.Done,
	)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
func (r *TaskRepository) FindAllForAdmin(ctx context.Context, filter dto.AdminTaskFilterQuery) ([]dto.AdminTaskResponse, int, error) {
	var conditions []string
	var args []interface{}

	if filter.UserID > 0 {
		conditions = append(conditions, "t.user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "t.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Priority != "" {
		conditions = append(conditions, "t.priority = ?")
		args = append(args, filter.Priority)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(t.title LIKE ? OR t.description LIKE ? OR u.name LIKE ? OR u.email LIKE ?)")
		kw := "%" + filter.Search + "%"
		args = append(args, kw, kw, kw, kw)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := `SELECT COUNT(1) FROM tasks t JOIN users u ON t.user_id = u.id ` + whereClause
	if err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	query := `SELECT t.id, t.user_id, u.name, u.email, t.category_id, c.name, t.title, t.description, t.priority, t.status, t.due_date, t.created_at, t.updated_at
	          FROM tasks t
	          JOIN users u ON t.user_id = u.id
	          LEFT JOIN categories c ON t.category_id = c.id ` + whereClause + `
	          ORDER BY t.created_at DESC LIMIT ? OFFSET ?`

	queryArgs := append(args, filter.Limit, offset)
	rows, err := r.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.AdminTaskResponse
	for rows.Next() {
		var item dto.AdminTaskResponse
		var catID sql.NullInt64
		var catName, desc sql.NullString
		var dueDate sql.NullTime
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserName, &item.UserEmail,
			&catID, &catName, &item.Title, &desc, &item.Priority, &item.Status,
			&dueDate, &createdAt, &updatedAt,
		); err != nil {
			return nil, 0, err
		}

		if catID.Valid {
			item.CategoryID = &catID.Int64
		}
		if catName.Valid {
			item.CategoryName = &catName.String
		}
		if desc.Valid {
			item.Description = &desc.String
		}
		if dueDate.Valid {
			formatted := dueDate.Time.Format("2006-01-02")
			item.DueDate = &formatted
		}
		if createdAt.Valid {
			item.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		if updatedAt.Valid {
			item.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	return items, total, rows.Err()
}