package models

import "golang.org/x/crypto/bcrypt"

func HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	User.Password = string(hashedPassword)
	return nil
}

func CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(password))
	return err == nil
}
