package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		gender TEXT NOT NULL,
		birth_date TIMESTAMPTZ NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user_friends (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		friend_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, friend_id),
		CONSTRAINT no_self_friend CHECK (user_id <> friend_id)
	);

	CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_gender ON users(gender);
	CREATE INDEX IF NOT EXISTS idx_users_birth_date ON users(birth_date);
	CREATE INDEX IF NOT EXISTS idx_user_friends_user ON user_friends(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_friends_friend ON user_friends(friend_id);
	`
	_, err := db.ExecContext(ctx, schema)
	return err
}

func SeedIfNeeded(ctx context.Context, db *sql.DB) error {
	var cnt int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&cnt); err != nil {
		return err
	}
	if cnt >= 20 {
		return nil
	}

	// Clean old data to have deterministic demo.
	_, _ = db.ExecContext(ctx, "DELETE FROM user_friends")
	_, _ = db.ExecContext(ctx, "DELETE FROM users")

	rnd := rand.New(rand.NewSource(42))
	genders := []string{"male", "female"}

	type seedUser struct {
		id        string
		name      string
		email     string
		gender    string
		birthDate time.Time
	}

	users := make([]seedUser, 0, 20)
	for i := 1; i <= 20; i++ {
		id := NewUUIDv4()
		name := fmt.Sprintf("User %02d", i)
		email := fmt.Sprintf("user%02d@mail.com", i)
		gender := genders[i%2]
		base := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
		days := rnd.Intn(365 * 9)
		birth := base.Add(time.Duration(days) * 24 * time.Hour)
		birth = time.Date(birth.Year(), birth.Month(), birth.Day(), i%24, 0, 0, 0, time.UTC)
		users = append(users, seedUser{id: id, name: name, email: email, gender: gender, birthDate: birth})
	}

	insertUser := `INSERT INTO users (id, name, email, gender, birth_date) VALUES ($1, $2, $3, $4, $5)`
	for _, u := range users {
		if _, err := db.ExecContext(ctx, insertUser, u.id, u.name, u.email, u.gender, u.birthDate); err != nil {
			return err
		}
	}

	// Make sure 2 users have >=3 common friends.
	u1 := users[0].id
	u2 := users[1].id
	common := []string{users[2].id, users[3].id, users[4].id}

	for _, f := range common {
		if err := insertFriend(ctx, db, u1, f); err != nil {
			return err
		}
		if err := insertFriend(ctx, db, u2, f); err != nil {
			return err
		}
	}

	// Add extra friends
	for _, f := range []string{users[5].id, users[6].id, users[7].id} {
		if err := insertFriend(ctx, db, u1, f); err != nil {
			return err
		}
	}
	for _, f := range []string{users[8].id, users[9].id} {
		if err := insertFriend(ctx, db, u2, f); err != nil {
			return err
		}
	}

	pairs := [][2]int{{3, 6}, {3, 7}, {4, 8}, {5, 9}, {6, 10}, {7, 11}, {8, 12}, {9, 13}, {10, 14}, {11, 15}, {12, 16}, {13, 17}, {14, 18}, {15, 19}, {16, 20}}
	for _, p := range pairs {
		uid := users[p[0]-1].id
		fid := users[p[1]-1].id
		if err := insertFriend(ctx, db, uid, fid); err != nil {
			return err
		}
	}

	return nil
}

func insertFriend(ctx context.Context, db *sql.DB, userID, friendID string) error {
	_, err := db.ExecContext(ctx, `INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, friendID)
	return err
}
