package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint64 `gorm:"primary_key"`
	Login    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
