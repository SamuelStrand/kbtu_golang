package repository

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type PaginatedUsers struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type UserFilters struct {
	ID        *string
	Name      string
	Email     string
	Gender    string
	BirthDate *time.Time
}

type UsersRepository interface {
	GetPaginatedUsers(limit, offset int, filters UserFilters, orderBy, orderDir string) (PaginatedUsers, error)
	GetCommonFriends(user1, user2 string) ([]User, error)
}
