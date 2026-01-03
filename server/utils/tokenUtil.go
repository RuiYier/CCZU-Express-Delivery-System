package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yurin-kami/PackChann/models"
)

type Claims struct {
	UserId    string
	StudentId string
	Role      string
	jwt.RegisteredClaims
}

func GenerateAllToken(user models.User, cfg models.JWTConfig) (string, string, error) {
	expirationTime := time.Now().Add(time.Duration(cfg.ExpirationHours) * time.Hour)

	claims := &Claims{
		UserId:    fmt.Sprint(user.UserId),
		StudentId: user.StudentId,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "PackChann",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Claims{
		UserId:    fmt.Sprint(user.UserId),
		StudentId: user.StudentId,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "PackChann",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}
