package model

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// User ユーザサービス構造体
type User struct {
	ID          uint32 `gorm:"primary_key"`
	Password    string `validate:"required_without=ID,max=32"`
	UserName    string `validate:"min=3,max=16"`
	Email       string `validate:"email"`
	Authority   uint32 `validate:"oneof=0 1 9"`
	ProfileText string `validate:"max=240"`
}
