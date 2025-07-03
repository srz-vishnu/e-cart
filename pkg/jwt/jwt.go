package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	jwtKey          = []byte("wA2I7VqLMbKP5RtUoD7M1jsJYWD9edxBS6cOgFXElwo=")
)

type Claims struct {
	UserID   int64  `json:"id"` //here we including id and name here so that will be there on the token, so we can use it in the other layers
	Username string `json:"username"`
	IsAdmin  bool   `json:"isadmin"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token
func GenerateToken(userID int64, username string, isadmin bool) (string, time.Time, error) {
	expirationTime := time.Now().Add(3 * time.Hour) // Token valid for 3 hours
	claims := &Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isadmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

// ValidateToken validates the JWT token and checks for expiration
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	// Check if the token has expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
