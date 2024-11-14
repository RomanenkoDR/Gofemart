package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/iternal/utils"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var reqUser utils.User

	err := json.NewDecoder(r.Body).Decode(&reqUser)

	if err != nil || reqUser.Login == "" || reqUser.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := getUserByLogin(reqUser.Login)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid login ", http.StatusUnauthorized)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !user.CheckPassword(reqUser.Password) {
		http.Error(w, "Invalid login ", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWT(user.Login)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusCreated)
}

func getUserByLogin(login string) (*utils.User, error) {
	var user utils.User
	err := database.QueryRow("SELECT id, login, password FROM users WHERE login = $1", login).Scan(
		&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
