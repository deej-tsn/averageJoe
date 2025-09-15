package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT_CustomClaim struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func CreateToken(username string, secretKey []byte) (string, error) {
	claims := &JWT_CustomClaim{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (*JWT_CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWT_CustomClaim{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if user, ok := token.Claims.(*JWT_CustomClaim); ok {
		return user, nil
	}
	return nil, errors.New("invalid token")
}
