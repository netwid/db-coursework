package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret")

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(id int) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)

	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JwtKey)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
