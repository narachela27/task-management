package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `INSERT INTO users (name, email, password, role, is_active, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, NOW(), NOW())`

	res, err := r.DB.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.IsActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at
	          FROM users WHERE email = ? LIMIT 1`

	var u model.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at
	          FROM users WHERE id = ? LIMIT 1`

	var u model.User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindAllUsers(ctx context.Context, filter dto.AdminUserFilterQuery) ([]model.User, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Role != "" {
		conditions = append(conditions, "role = ?")
		args = append(args, filter.Role)
	}
	if filter.IsActive == "true" || filter.IsActive == "1" {
		conditions = append(conditions, "is_active = 1")
	} else if filter.IsActive == "false" || filter.IsActive == "0" {
		conditions = append(conditions, "is_active = 0")
	}
	if filter.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR email LIKE ?)")
		kw := "%" + filter.Search + "%"
		args = append(args, kw, kw)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM users "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at
	          FROM users ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	queryArgs := append(args, filter.Limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepository) UpdateRole(ctx context.Context, id int64, role string) (bool, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE users SET role = ?, updated_at = NOW() WHERE id = ?`, role, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, isActive bool) (bool, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE users SET is_active = ?, updated_at = NOW() WHERE id = ?`, isActive, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}