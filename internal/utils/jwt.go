package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
var refreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

type JWTClaims struct {
	IDUser   string  `json:"id_user"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

type CreateJWTPayload struct {
	IDUser   string
	Username string
	Email    string
}

func CreateAccessToken(payload *CreateJWTPayload) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &JWTClaims{
		IDUser:   payload.IDUser,
		Username: &payload.Username,
		Email:    &payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rizz",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(accessTokenSecret)
}

func CreateRefreshToken(idUser string) (string, error) {
	expirationTime := time.Now().Add(7 * time.Hour)

	claims := &JWTClaims{
		IDUser: idUser, // * Hanya butuh ID untuk refresh
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rizz-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(refreshTokenSecret)
}

func ValidateToken(tokenString string, secretKey []byte) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
