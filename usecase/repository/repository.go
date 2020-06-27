package repository

import "github.com/yzmw1213/GoMicroApp/domain/model"

// BlogRepository 投稿サービスの抽象定義
type BlogRepository interface {
	Create(*model.Blog) error
	Delete(*model.Blog) error
	List() ([]model.Blog, error)
	Update(*model.Blog) error
}
