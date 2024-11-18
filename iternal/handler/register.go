package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	jwtKey   []byte
	database *gorm.DB
)

func initJwtKey() {
	key, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		log.Fatal("No secret key provided")
	}
	jwtKey = []byte(key)
}

func SetDatabase(db *gorm.DB) {
	database = db
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Login == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя
	var existingUser models.User
	result := database.Where("login = ?", user.Login).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Хэшируем пароль (предполагается, что метод `HashPassword` уже существует)
	err = user.HashPassword()
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Сохранение пользователя
	err = database.Create(&user).Error
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Генерация токена
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
