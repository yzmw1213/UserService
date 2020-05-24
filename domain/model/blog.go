package model

import "github.com/jinzhu/gorm"

type Blog struct {
	gorm.Model
	AuthorId int32
	Title    string
	Content  string
}
