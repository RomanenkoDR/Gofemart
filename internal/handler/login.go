package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services"
	"log"
	"net/http"
	"time"
)

// Login обрабатывает аутентификацию пользователя.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.User

	// Декодируем данные из запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Login == "" || input.Password == "" {
		log.Printf("Некорректный запрос \n%v", err)
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы
	storedUser := &models.User{}
	if err := db.GetUserByLogin(h.DB, input.Login, storedUser); err != nil {
		log.Printf("Пользователь не авторизован \n%v", err)
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !services.CheckPassword(input.Password, storedUser.Password) {
		log.Printf("Неверные логин и/или пароль \n%v", http.StatusUnauthorized)
		http.Error(w, "неверные логин и/или пароль", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	jwtToken, err := services.GenerateJWT(storedUser.Login)
	if err != nil {
		log.Printf("ошибка при генерации JWT-токе \n%v", err)
		http.Error(w, "ошибка при генерации JWT-токена", http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Authorization", "Bearer "+jwtToken)
	w.WriteHeader(http.StatusOK)
}
