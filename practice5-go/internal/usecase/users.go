package usecase

import (
	"context"
	"fmt"
	"strings"

	"practice5/internal/repository"
)

type UserFilters struct {
	ID        string
	Name      string
	Email     string
	Gender    string
	BirthDate string // YYYY-MM-DD
}

type UsersUsecase struct {
	repo repository.UsersRepository
}

func NewUsersUsecase(repo repository.UsersRepository) *UsersUsecase {
	return &UsersUsecase{repo: repo}
}

func (u *UsersUsecase) GetPaginatedUsers(ctx context.Context, limit, offset int, f UserFilters, orderBy, orderDir string) (repository.PaginatedUsers, error) {
	_ = ctx

	filters := repository.UserFilters{
		Name:   strings.TrimSpace(f.Name),
		Email:  strings.TrimSpace(f.Email),
		Gender: strings.TrimSpace(f.Gender),
	}

	if strings.TrimSpace(f.ID) != "" {
		id := strings.TrimSpace(f.ID)
		filters.ID = &id
	}

	if strings.TrimSpace(f.BirthDate) != "" {
		bd, err := repository.ParseBirthDate(f.BirthDate)
		if err != nil {
			return repository.PaginatedUsers{}, fmt.Errorf("birth_date must be YYYY-MM-DD")
		}
		filters.BirthDate = bd
	}

	return u.repo.GetPaginatedUsers(limit, offset, filters, orderBy, orderDir)
}

func (u *UsersUsecase) GetCommonFriends(ctx context.Context, user1, user2 string) ([]repository.User, error) {
	_ = ctx
	return u.repo.GetCommonFriends(user1, user2)
}
