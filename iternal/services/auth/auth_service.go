package services

import (
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	Claims       jwt.MapClaims
	User         models.User
	JwtKey       string
	Database     *gorm.DB
	ExistingUser models.User
)
