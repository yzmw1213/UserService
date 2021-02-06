package repository

import "github.com/yzmw1213/UserService/domain/model"

// UserRepository ユーザーサービスの抽象定義
type UserRepository interface {
	LoginAuth(string, string) (*model.Auth, error)
	Create(*model.User) (*model.User, error)
	CreateDemoUser() (*model.Auth, error)
	GetUserByEmail(string) (model.User, error)
	OtherUserExistsByEmail(string, uint32) bool
	DeleteByID(id uint32) error
	Follow(*model.Relation) (*model.Relation, error)
	UnFollow(*model.Relation) (*model.Relation, error)
	List() ([]model.User, error)
	ListAllNormalUser() ([]model.User, error)
	Update(*model.User) (*model.User, error)
	GetFollowUsersByID(id uint32) []uint32
}
