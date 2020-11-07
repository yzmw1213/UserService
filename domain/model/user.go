package model

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// User ユーザサービス構造体
type User struct {
	ID          uint32 `gorm:"primary_key"`
	Password    string `validate:"min=6,max=32"`
	UserName    string `validate:"min=6,max=16"`
	Email       string `validate:"email"`
	Gender      uint32 `validate:"oneof=0 1 2 9"`
	Authority   uint32 `validate:"oneof=0 1 2 9"`
	ProfileText string `validate:"max=240"`
}
