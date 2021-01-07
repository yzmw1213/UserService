package repository

import "github.com/yzmw1213/UserService/domain/model"

// UserRepository 投稿サービスの抽象定義
type UserRepository interface {
	LoginAuth(string, string) (*model.Auth, error)
	Create(*model.User) (*model.User, error)
	Read(uint32) (model.User, error)
	GetUserByEmail(string) (model.User, error)
	OtherUserExistsByEmail(string, uint32) bool
	DeleteByID(id uint32) error
	List() ([]model.User, error)
	ListAllNormalUser() ([]model.User, error)
	Update(*model.User) (*model.User, error)
}
