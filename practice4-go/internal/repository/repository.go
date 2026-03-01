package repository

import (
	"context"

	"practice4/internal/repository/_postgres"
	"practice4/internal/repository/_postgres/users"
	"practice4/pkg/modules"
)

type UserRepository interface {
	GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, user modules.User) (int, error)
	UpdateUser(ctx context.Context, id int, user modules.User) error
	DeleteUser(ctx context.Context, id int) (int64, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
