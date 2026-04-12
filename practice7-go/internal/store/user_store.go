package store

import (
	"errors"
	"sync"

	"practice7/internal/entity"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrUsernameTaken = errors.New("username already exists")
	ErrEmailTaken    = errors.New("email already exists")
	ErrInvalidRole   = errors.New("invalid role")
)

type UserStore interface {
	Create(u *entity.User) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetByID(id uuid.UUID) (*entity.User, error)
	UpdateRole(id uuid.UUID, role string) (*entity.User, error)
}

type InMemoryUserStore struct {
	mu         sync.RWMutex
	byID       map[uuid.UUID]*entity.User
	byUsername map[string]uuid.UUID
	byEmail    map[string]uuid.UUID
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		byID:       make(map[uuid.UUID]*entity.User),
		byUsername: make(map[string]uuid.UUID),
		byEmail:    make(map[string]uuid.UUID),
	}
}

func (s *InMemoryUserStore) Create(u *entity.User) (*entity.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byUsername[u.Username]; ok {
		return nil, ErrUsernameTaken
	}
	if _, ok := s.byEmail[u.Email]; ok {
		return nil, ErrEmailTaken
	}

	clone := *u
	s.byID[clone.ID] = &clone
	s.byUsername[clone.Username] = clone.ID
	s.byEmail[clone.Email] = clone.ID

	return publicUser(&clone), nil
}

func (s *InMemoryUserStore) GetByUsername(username string) (*entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byUsername[username]
	if !ok {
		return nil, ErrNotFound
	}
	u := s.byID[id]
	clone := *u
	return &clone, nil
}

func (s *InMemoryUserStore) GetByID(id uuid.UUID) (*entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return publicUser(u), nil
}

func (s *InMemoryUserStore) UpdateRole(id uuid.UUID, role string) (*entity.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	u.Role = role
	return publicUser(u), nil
}

func publicUser(u *entity.User) *entity.User {
	clone := *u
	clone.PasswordHash = ""
	return &clone
}
