package model

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// User ユーザサービス構造体
type User struct {
	UserID   int32  `gorm:"primary_key"`
	Password string `validate:"min=6,max=32"`
	UserName string `validate:"min=6,max=16"`
	Email    string `validate:"email"`
}
