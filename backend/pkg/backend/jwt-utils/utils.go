package jwt_utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func GenerateFromLogin(login string, jwtKey string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
