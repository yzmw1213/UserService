package repository

import "github.com/yzmw1213/GoMicroApp/domain/model"

type BlogRepository interface {
	Create(*model.Blog) error
	Delete(*model.Blog) error
}
