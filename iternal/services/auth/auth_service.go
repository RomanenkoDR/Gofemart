package auth

import (
	"errors"
	"fmt"
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

var (
	Claims       jwt.MapClaims
	User         models.User
	JwtKey       string
	ExistingUser models.User
)

func СheckAuthToken(r *http.Request) (string, int, error) {
	jwtToken := r.Header.Get("Authorization")

	if jwtToken == "" {
		return "", http.StatusUnauthorized, errors.New("пользователь не авторизован")
	}

	t, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	switch {
	case t.Valid:
		Claims = t.Claims.(jwt.MapClaims)
		username, ok := Claims["username"].(string)
		if !ok || username == "" {
			return "", http.StatusBadRequest, errors.New("неудалось получить username из токена JWT")
		}
		return username, http.StatusOK, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return "", http.StatusBadRequest, fmt.Errorf("токен имеет неправильную форму %w", err)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return "", http.StatusUnauthorized, fmt.Errorf("подпись токена недействительна %w", err)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return "", http.StatusUnauthorized, fmt.Errorf("срок действия токена истек или токен еще не действителен %w", err)
	default:
		return "", http.StatusInternalServerError, fmt.Errorf("неизвестная ошика на этапе проверки валидности токена: %w", err)
	}
}
