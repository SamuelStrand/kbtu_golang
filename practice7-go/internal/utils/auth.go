package utils

import (
	"errors"
	"fmt"
	"strings"

	"practice7/internal/entity"
	"practice7/internal/store"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthService struct {
	store  store.UserStore
	secret []byte
}

func NewAuthService(s store.UserStore, secret []byte) *AuthService {
	return &AuthService{store: s, secret: secret}
}

func (a *AuthService) Register(req RegisterRequest) (*entity.User, error) {
	role := "user"
	if strings.TrimSpace(req.Role) != "" {
		if req.Role != "user" && req.Role != "admin" {
			return nil, store.ErrInvalidRole
		}
		role = req.Role
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &entity.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         role,
	}

	return a.store.Create(u)
}

func (a *AuthService) Login(req LoginRequest) (string, *entity.User, error) {
	u, err := a.store.GetByUsername(req.Username)
	if err != nil || !CheckPassword(u.PasswordHash, req.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := GenerateJWT(u.ID, u.Role, a.secret)
	if err != nil {
		return "", nil, fmt.Errorf("generate jwt: %w", err)
	}

	public := &entity.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
	return token, public, nil
}

func (a *AuthService) GetByID(id uuid.UUID) (*entity.User, error) {
	return a.store.GetByID(id)
}

func (a *AuthService) PromoteUserToAdmin(id uuid.UUID) (*entity.User, error) {
	return a.store.UpdateRole(id, "admin")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
