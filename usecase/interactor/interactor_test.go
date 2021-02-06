package interactor

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/UserService/db"
	"github.com/yzmw1213/UserService/domain/model"
)

var (
	testemail       = "test@example.com"
	superemail      = "super@example.com"
	testpassword    = "password"
	demoProfileText = "プロフィールが入ります"
	updatedName     = "updatedName"
	user1           uint32
	user2           uint32
	user3           uint32
)

var DemoUser = model.User{
	UserName:    "testuser",
	Password:    testpassword,
	Email:       testemail,
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var DemoSuperUser = model.User{
	UserName:  "superuser",
	Password:  testpassword,
	Email:     superemail,
	Authority: authoritySuperUser,
}

var NameNullUser = &model.User{
	UserName:    "",
	Password:    testpassword,
	Email:       testemail,
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var NameTooLongUser = &model.User{
	UserName:    "testusertestusert",
	Password:    testpassword,
	Email:       testemail,
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var NameTooShortUser = &model.User{
	UserName:    "testu",
	Password:    testpassword,
	Email:       testemail,
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var EmailNullUser = &model.User{
	UserName:    "testuser",
	Password:    testpassword,
	Email:       "",
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var EmailInvalidUser1 = &model.User{
	UserName:    "testuser",
	Password:    testpassword,
	Email:       "test",
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

var EmailInvalidUser2 = &model.User{
	UserName:    "testuser",
	Password:    testpassword,
	Email:       "@example.com",
	Authority:   authorityNormalUser,
	ProfileText: demoProfileText,
}

// TestCreate ユーザー作成の正常系
func TestCreate(t *testing.T) {
	initUserTable()
	var i UserInteractor
	var users []model.User
	// user := &DemoUser
	users = append(users, DemoUser)
	DemoUser.Email = "acregh@example.com"
	users = append(users, DemoUser)
	DemoUser.Email = "uiwefg@example.com"
	users = append(users, DemoUser)
	DemoUser.Email = "dgjs5ts@example.com"
	users = append(users, DemoUser)
	DemoUser.Email = "th5fa6j@example.com"
	users = append(users, DemoUser)
	for _, u := range users {
		createdUser, err := i.Create(&u)
		assert.Equal(t, nil, err)
		assert.Equal(t, u.UserName, createdUser.UserName)
		assert.Equal(t, u.Email, createdUser.Email)
		assert.Equal(t, u.Password, createdUser.Password)
		assert.Equal(t, authorityNormalUser, createdUser.Authority)
	}
}

func TestCount(t *testing.T) {
	var i UserInteractor
	var u model.User
	count, err := i.Count(u)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, count)

}

// TestCreateSuperUser 管理者ユーザー作成の正常系
func TestCreateSuperUser(t *testing.T) {
	var i UserInteractor
	user := &DemoSuperUser
	createdUser, err := i.Create(user)

	assert.Equal(t, nil, err)
	assert.Equal(t, user.UserName, createdUser.UserName)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
	assert.Equal(t, authoritySuperUser, createdUser.Authority)
}

func TestCreateUserEmailUsed(t *testing.T) {
	var i UserInteractor
	user := &DemoUser
	_, err := i.Create(user)

	assert.NotEqual(t, nil, err)
}

func TestGetUserByEmail(t *testing.T) {
	var i UserInteractor
	email := DemoUser.Email
	user, err := i.GetUserByEmail(email)

	assert.Equal(t, nil, err)
	assert.Equal(t, DemoUser.UserName, user.UserName)
	assert.Equal(t, DemoUser.Email, user.Email)
}

func TestCreateNameNull(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(NameNullUser)

	assert.NotEqual(t, nil, err)
}

func TestCreateNameTooLong(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(NameTooLongUser)

	assert.NotEqual(t, nil, err)
}

func TestCreateNameTooShort(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(NameTooShortUser)

	assert.NotEqual(t, nil, err)
}

func TestCreateEmailNull(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(EmailNullUser)

	assert.NotEqual(t, nil, err)
}

func TestCreateEmailInvalid1(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(EmailInvalidUser1)

	assert.NotEqual(t, nil, err)
}

func TestCreateEmailInvalid2(t *testing.T) {
	var i UserInteractor
	_, err := i.Create(EmailInvalidUser2)

	assert.NotEqual(t, nil, err)
}

func TestUpdate(t *testing.T) {
	var i UserInteractor

	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)
	inputUser := findUser

	inputUser.UserName = updatedName
	inputUser.Password = "password"

	updatedUser, err := i.Update(&inputUser)

	assert.Equal(t, nil, err)
	assert.Equal(t, updatedUser.ID, findUser.ID)
	assert.Equal(t, updatedUser.Email, findUser.Email)
	assert.NotEqual(t, updatedUser.Password, findUser.Password)
	assert.NotEqual(t, updatedUser.UserName, findUser.UserName)

	// err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(inputUser.Password))
	// assert.Equal(t, nil, err)

}

func TestUpdateNameNull(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)
	inputUser := findUser

	inputUser.UserName = ""
	inputUser.Password = "password"

	_, err = i.Update(&inputUser)
	assert.NotEqual(t, nil, err)
}

func TestUpdateNameTooLong(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)
	inputUser := findUser

	inputUser.UserName = "testusertestusert"
	inputUser.Password = "password"

	_, err = i.Update(&inputUser)
	assert.NotEqual(t, nil, err)
}

func TestUpdateNameTooShort(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)
	inputUser := findUser

	inputUser.UserName = "testu"
	inputUser.Password = "password"

	_, err = i.Update(&inputUser)
	assert.NotEqual(t, nil, err)
}

func TestUpdateUserEmailUsed(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)
	inputUser := findUser

	inputUser.Email = testemail
	inputUser.Password = "password"

	_, err = i.Update(&inputUser)
	assert.NotEqual(t, nil, err)
}

func TestRead(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)

	user, err := i.Read(findUser.ID)
	assert.Equal(t, nil, err)
	assert.Equal(t, updatedName, user.UserName)
	assert.Equal(t, DemoUser.Email, user.Email)

}

func TestReadByIDNotExists(t *testing.T) {
	var i UserInteractor
	var searchID uint32 = 10000

	_, err := i.Read(searchID)
	assert.NotEqual(t, nil, err)

}

func TestLoginAuth(t *testing.T) {
	var i UserInteractor
	email := testemail
	password := testpassword

	auth, err := i.LoginAuth(email, password)

	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", auth.Token)
	assert.NotEqual(t, zero, auth.UserID)
}

// TestLoginAuthPasswordInvalid 登録したパスワードと異なるパスワードでログインを行う異常系
func TestLoginAuthPasswordInvalid(t *testing.T) {
	var i UserInteractor
	email := testemail
	password := "aaaaaa"

	auth, err := i.LoginAuth(email, password)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", auth.Token)

}

func TestFollow(t *testing.T) {
	var i UserInteractor
	user1, user2, user3 := selectUsers()
	assert.Equal(t, nil, err)

	// postID := createdPost.Post.ID
	relation := &model.Relation{FollowerUserID: user1, FollowedUserID: user2}

	_, err = i.Follow(relation)
	assert.Equal(t, nil, err)

	// フォロー数をカウントするテスト
	likeCount := countFollowUserByFollower(user1)
	assert.Equal(t, 1, likeCount)

	relation = &model.Relation{FollowerUserID: user1, FollowedUserID: user3}
	_, err = i.Follow(relation)

	// フォロー数が1増えている事をテスト
	likeCount = countFollowUserByFollower(user1)
	assert.Equal(t, 2, likeCount)
}

func TestUnFollow(t *testing.T) {
	var i UserInteractor

	// フォロワー・フォローユーザー関係
	relations := []model.Relation{
		{FollowerUserID: user2, FollowedUserID: user1},
		{FollowerUserID: user2, FollowedUserID: user3},
	}
	// フォロー実行
	for _, r := range relations {
		_, err = i.Follow(&r)
	}
	assert.Equal(t, nil, err)

	beforeFollowCount := countFollowUserByFollower(user2)

	// 1件フォロー解除
	_, err = i.UnFollow(&model.Relation{FollowerUserID: user2, FollowedUserID: user1})
	assert.Equal(t, nil, err)

	afterFollowCount := countFollowUserByFollower(user2)

	// フォロー数が1だけ減っている事をテスト
	assert.Equal(t, afterFollowCount, beforeFollowCount-1)
}

// TestLoginAuthPasswordNull パスワード空白でログインを行う異常系
func TestLoginAuthPasswordNull(t *testing.T) {
	var i UserInteractor
	email := testemail
	password := ""

	auth, err := i.LoginAuth(email, password)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", auth.Token)

}

// TestLoginAuthInvalidEmail 登録のないEmailでログインを行う異常系
func TestLoginAuthInvalidEmail(t *testing.T) {
	var i UserInteractor
	email := "notused@gmail.com"
	password := testpassword

	auth, err := i.LoginAuth(email, password)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", auth.Token)
}

// TestLoginAuthEmailNull Email入力空白でログインを行う異常系
func TestLoginAuthEmailNull(t *testing.T) {
	var i UserInteractor
	email := ""
	password := testpassword

	auth, err := i.LoginAuth(email, password)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", auth.Token)
}

// TestDeleteNotExistsUser 登録のないユーザを削除する異常系
func TestDeleteNotExistsUser(t *testing.T) {
	var i UserInteractor
	var deleteUser *model.User

	var u model.User
	countBeforeDelete, err := i.Count(u)
	assert.Equal(t, nil, err)

	findUser, err := i.GetUserByEmail("notused@gmail.com")
	assert.NotEqual(t, nil, err)

	deleteUser = &findUser
	deleteUser.ID = 10000

	err = i.DeleteByID(10000)
	assert.Equal(t, nil, err)

	countAfterDelete, err := i.Count(u)
	assert.Equal(t, nil, err)
	assert.Equal(t, countBeforeDelete, countAfterDelete)
}

func TestDelete(t *testing.T) {
	var i UserInteractor
	var deleteUser *model.User
	findUser, err := i.GetUserByEmail(testemail)
	assert.Equal(t, nil, err)

	deleteUser = &findUser

	err = i.DeleteByID(deleteUser.ID)
	assert.Equal(t, nil, err)

	findUser, err = i.GetUserByEmail(testemail)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, zero, findUser.ID)
	assert.Equal(t, "", findUser.UserName)
	assert.Equal(t, "", findUser.Email)
}

func selectUsers() (uint32, uint32, uint32) {
	var i UserInteractor
	var num int = 1
	users, _ := i.ListAllNormalUser()

	for i := range users {
		if num == 1 {
			user1 = users[i].ID
			log.Println("user1", user1)
		}
		if num == 2 {
			user2 = users[i].ID
			log.Println("user2", user2)
		}
		if num == 3 {
			user3 = users[i].ID
			log.Println("user3", user3)
		}
		num++
	}
	return user1, user2, user3
}

func initUserTable() {
	DB := db.GetDB()
	DB.Delete(&model.User{})
	DB.Delete(&model.Relation{})
}
