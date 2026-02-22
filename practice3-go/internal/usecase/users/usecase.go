package users

import (
	"context"
	"errors"
	"strings"

	"practice3/internal/repository"
	repoUsers "practice3/internal/repository/_postgres/users"
	"practice3/pkg/modules"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrBadRequest  = errors.New("bad request")
	ErrEmailFormat = errors.New("invalid email")
)

type Usecase struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error) {
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.GetUsers(ctx, limit, offset)
}

func (u *Usecase) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	if id <= 0 {
		return nil, ErrBadRequest
	}
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoUsers.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *Usecase) CreateUser(ctx context.Context, user modules.User) (int, error) {
	if strings.TrimSpace(user.Name) == "" {
		return 0, ErrBadRequest
	}
	if !strings.Contains(user.Email, "@") {
		return 0, ErrEmailFormat
	}
	if user.Age < 0 || user.Age > 150 {
		return 0, ErrBadRequest
	}
	return u.repo.CreateUser(ctx, user)
}

func (u *Usecase) UpdateUser(ctx context.Context, id int, user modules.User) error {
	if id <= 0 {
		return ErrBadRequest
	}
	if strings.TrimSpace(user.Name) == "" {
		return ErrBadRequest
	}
	if !strings.Contains(user.Email, "@") {
		return ErrEmailFormat
	}
	if user.Age < 0 || user.Age > 150 {
		return ErrBadRequest
	}

	if err := u.repo.UpdateUser(ctx, id, user); err != nil {
		if errors.Is(err, repoUsers.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (u *Usecase) DeleteUser(ctx context.Context, id int) (int64, error) {
	if id <= 0 {
		return 0, ErrBadRequest
	}
	affected, err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repoUsers.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return affected, nil
}
