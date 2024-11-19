package config

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Options struct {
	Host string `env:"SERV_HOST" envDefault:"localhost"`
	Port string `env:"SERV_PORT" envDefault:"8080"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {

	// Получаем SECRET_KEY из переменной окружения
	jwtKey := os.Getenv("SECRET_KEY")
	if jwtKey == "" {
		return "", fmt.Errorf("SECRET_KEY is not set in environment variables")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
