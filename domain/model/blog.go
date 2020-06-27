package model

import "time"

// Blog 投稿サービス構造体
type Blog struct {
	BlogID    int32 `gorm:"primary_key"`
	AuthorID  int32
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
