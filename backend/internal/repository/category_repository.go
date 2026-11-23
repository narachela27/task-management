package repository

import (
	"context"
	"database/sql"
	"errors"

	"task-management-backend/internal/model"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Create(ctx context.Context, cat *model.Category) (int64, error) {
	query := `INSERT INTO categories (user_id, name, created_at, updated_at)
	          VALUES (?, ?, NOW(), NOW())`

	res, err := r.DB.ExecContext(ctx, query, cat.UserID, cat.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CategoryRepository) FindAllByUserID(ctx context.Context, userID int64) ([]model.Category, error) {
	query := `SELECT id, user_id, name, created_at, updated_at
	          FROM categories WHERE user_id = ? ORDER BY name ASC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*model.Category, error) {
	query := `SELECT id, user_id, name, created_at, updated_at
	          FROM categories WHERE id = ? AND user_id = ? LIMIT 1`

	var c model.Category
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&c.ID, &c.UserID, &c.Name, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) ExistsByNameAndUserID(ctx context.Context, name string, userID int64, excludeID int64) (bool, error) {
	query := `SELECT COUNT(1) FROM categories WHERE user_id = ? AND LOWER(name) = LOWER(?) AND id != ?`
	var count int
	err := r.DB.QueryRowContext(ctx, query, userID, name, excludeID).Scan(&count)
	return count > 0, err
}

func (r *CategoryRepository) Update(ctx context.Context, id, userID int64, name string) (bool, error) {
	query := `UPDATE categories SET name = ?, updated_at = NOW() WHERE id = ? AND user_id = ?`
	res, err := r.DB.ExecContext(ctx, query, name, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (r *CategoryRepository) Delete(ctx context.Context, id, userID int64) (bool, error) {
	query := `DELETE FROM categories WHERE id = ? AND user_id = ?`
	res, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}
