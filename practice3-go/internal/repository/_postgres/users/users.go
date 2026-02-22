package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"practice3/internal/repository/_postgres"
	"practice3/pkg/modules"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var users []modules.User
	q := `SELECT id, name, email, age, created_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2`
	if err := r.db.DB.SelectContext(ctx, &users, q, limit, offset); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var u modules.User
	q := `SELECT id, name, email, age, created_at, deleted_at
		FROM users
		WHERE id=$1 AND deleted_at IS NULL`
	if err := r.db.DB.GetContext(ctx, &u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateUser(ctx context.Context, user modules.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var newID int
	q := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	if err := tx.QueryRowContext(ctx, q, user.Name, user.Email, user.Age).Scan(&newID); err != nil {
		return 0, err
	}

	q2 := `INSERT INTO audit_logs (user_id, action) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, q2, newID, "CREATE_USER"); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, user modules.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	q := `UPDATE users
		SET name=$1, email=$2, age=$3
		WHERE id=$4 AND deleted_at IS NULL`
	res, err := r.db.DB.ExecContext(ctx, q, user.Name, user.Email, user.Age, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	q := `UPDATE users SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
	res, err := r.db.DB.ExecContext(ctx, q, id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		return 0, ErrNotFound
	}

	q2 := `INSERT INTO audit_logs (user_id, action) VALUES ($1, $2)`
	_, _ = r.db.DB.ExecContext(ctx, q2, id, "SOFT_DELETE_USER")

	return affected, nil
}
