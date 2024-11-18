package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil || reqUser.Login == "" || reqUser.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User

	result := database.Where("login = ?", reqUser.Login).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if result.Error != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !user.CheckPassword(reqUser.Password) {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
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
