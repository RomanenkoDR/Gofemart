package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/iternal/utils"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	jwtKey   []byte
	database *sql.DB
)

func initJwtKey() {
	key, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		log.Fatal("No secret key provided")
	}
	jwtKey = []byte(key)
}

func SetDatabase(db *sql.DB) {
	database = db // Устанавливаем глобальную переменную
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user utils.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Login == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	exists, err := isLoginExists(user.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already taken", http.StatusConflict)
		return
	}

	err = saveUserToDB(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tokenString, err := generateJWT(user.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}

func isLoginExists(login string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)"
	err := database.QueryRow(query, login).Scan(&exists)
	return exists, err
}

func saveUserToDB(user utils.User) error {
	_, err := database.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", user.Login, user.Password)
	return err

}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
