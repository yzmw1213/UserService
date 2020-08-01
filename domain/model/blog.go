package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Blog 投稿サービス構造体
type Blog struct {
	BlogID    int32  `gorm:"primary_key"`
	AuthorID  int32  `validate:"required,number"`
	Title     string `validate:"min=1,max=32"`
	Content   string `validate:"min=1,max=140"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
