package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Balance float64 `json:"balance"`
	jwt.RegisteredClaims
}

func GetToken(balance float64) (string, error) {
	claims := Claims{
		Balance: balance,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "employee_test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret_code"))
}

func CheckToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_code"), nil
	})

	if err != nil || !token.Valid {
		return "Invalid token : ", err
	}

	return "Valid token", nil
}
