package auth

import (
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
