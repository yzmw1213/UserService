package model

import "time"

type Blog struct {
	BlogId    int32 `gorm:"primary_key"`
	AuthorId  int32
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
