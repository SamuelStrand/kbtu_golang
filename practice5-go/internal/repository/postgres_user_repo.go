package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

var allowedOrderColumns = map[string]string{
	"id":         "u.id",
	"name":       "u.name",
	"email":      "u.email",
	"gender":     "u.gender",
	"birth_date": "u.birth_date",
}

func (r *PostgresUserRepo) GetPaginatedUsers(limit, offset int, filters UserFilters, orderBy, orderDir string) (PaginatedUsers, error) {
	whereParts := make([]string, 0, 5)
	args := make([]any, 0, 7)
	argN := 1

	if filters.ID != nil {
		whereParts = append(whereParts, fmt.Sprintf("u.id = $%d", argN))
		args = append(args, *filters.ID)
		argN++
	}
	if strings.TrimSpace(filters.Name) != "" {
		whereParts = append(whereParts, fmt.Sprintf("u.name ILIKE $%d", argN))
		args = append(args, "%"+filters.Name+"%")
		argN++
	}
	if strings.TrimSpace(filters.Email) != "" {
		whereParts = append(whereParts, fmt.Sprintf("u.email ILIKE $%d", argN))
		args = append(args, "%"+filters.Email+"%")
		argN++
	}
	if strings.TrimSpace(filters.Gender) != "" {
		whereParts = append(whereParts, fmt.Sprintf("u.gender = $%d", argN))
		args = append(args, filters.Gender)
		argN++
	}
	if filters.BirthDate != nil {
		whereParts = append(whereParts, fmt.Sprintf("u.birth_date::date = $%d", argN))
		args = append(args, filters.BirthDate.Format("2006-01-02"))
		argN++
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	orderCol, ok := allowedOrderColumns[strings.ToLower(strings.TrimSpace(orderBy))]
	if !ok {
		orderCol = "u.id" // default if nothing provided
	}
	dir := strings.ToUpper(strings.TrimSpace(orderDir))
	if dir != "ASC" && dir != "DESC" {
		dir = "ASC"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u %s", whereSQL)
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return PaginatedUsers{}, err
	}

	query := fmt.Sprintf(
		`SELECT u.id, u.name, u.email, u.gender, u.birth_date
		 FROM users u
		 %s
		 ORDER BY %s %s
		 LIMIT $%d OFFSET $%d`,
		whereSQL, orderCol, dir, argN, argN+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return PaginatedUsers{}, err
	}
	defer rows.Close()

	users := make([]User, 0, limit)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return PaginatedUsers{}, err
		}
		users = append(users, u)
	}

	return PaginatedUsers{Data: users, TotalCount: total, Limit: limit, Offset: offset}, nil
}

func (r *PostgresUserRepo) GetCommonFriends(user1, user2 string) ([]User, error) {
	if user1 == user2 {
		return nil, errors.New("user1 and user2 must be different")
	}

	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
		ORDER BY u.name ASC
	`

	rows, err := r.db.Query(query, user1, user2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := make([]User, 0, 10)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func ParseBirthDate(dateStr string) (*time.Time, error) {
	if strings.TrimSpace(dateStr) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
