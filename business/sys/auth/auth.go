package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(tkn string) (Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tkn, &claims, keyFunc)
	if err != nil {
		return Claims{}, fmt.Errorf("parsing token failed: %w", err)
	}

	if !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing token failed: %w", err)
	}

	return tkn, nil
}

// private

// TODO: extract it to ENV or somewhere else.
const secret = "secret"

func keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(secret), nil
}
