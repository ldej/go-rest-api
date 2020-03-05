package jwt

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserUID string `json:"uid"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

func CreateToken(jwtKey []byte, uid string, role string, expirationTime time.Time) (string, error) {
	claims := Claims{
		UserUID:        uid,
		Role:           role,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	return token.SignedString(jwtKey)
}

func ValidateToken(jwtKey []byte, tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, fmt.Errorf("invalid signature %w", err)
		}
		return claims, err
	}
	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}
	return claims, nil
}
