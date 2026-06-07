package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId   int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var secret_key = []byte("SecretKey")

func GenerateToken(user_id int64, username string) (string, error) {
	claims := Claims{
		UserId:   user_id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret_key)
}

func ValidateToken(token_string string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		token_string,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret_key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
