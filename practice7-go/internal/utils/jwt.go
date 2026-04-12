package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string
	Role   string
}

func GenerateJWT(userID uuid.UUID, role string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseJWT(tokenStr string, secret []byte) (JWTClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return JWTClaims{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JWTClaims{}, errors.New("invalid claims")
	}
	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)
	if userID == "" || role == "" {
		return JWTClaims{}, errors.New("missing claims")
	}
	return JWTClaims{UserID: userID, Role: role}, nil
}

func ParseJWTLoose(tokenStr string, secret []byte) (JWTClaims, error) {
	return ParseJWT(tokenStr, secret)
}
