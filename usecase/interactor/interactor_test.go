package interactor

import (
	"testing"

	"github.com/yzmw1213/UserService/db"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/UserService/domain/model"
)

var (
	testemail    = "test@gmail.com"
	testpassword = "password"
	updatedName  = "updatedName"
)

var DemoUser = model.User{
	UserName: "testuser",
	Password: testpassword,
	Email:    testemail,
}

var NameNullUser = &model.User{
	UserName: "",
	Password: testpassword,
	Email:    testemail,
}

var NameTooLongUser = &model.User{
	UserName: "testusertestusert",
	Password: testpassword,
	Email:    testemail,
}

var NameTooShortUser = &model.User{
	UserName: "testu",
	Password: testpassword,
	Email:    testemail,
}

var EmailNullUser = &model.User{
	UserName: "testuser",
	Password: testpassword,
	Email:    "",
}

var EmailInvalidUser1 = &model.User{
	UserName: "testuser",
	Password: testpassword,
	Email:    "test",
}

var EmailInvalidUser2 = &model.User{
	UserName: "testuser",
	Password: testpassword,
	Email:    "@gmail.com",
}

// TestCreate ユーザー作成の正常系
func TestCreate(t *testing.T) {
	initUserTable()
	var i UserInteractor
	user := &DemoUser
	createdUser, err := i.Create(user)

	assert.Equal(t, nil, err)
	assert.Equal(t, user.UserName, createdUser.UserName)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
}

func TestCount(t *testing.T) {
	var i UserInteractor
	var u model.User
	count, err := i.Count(u)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, count)

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
	assert.Equal(t, DemoUser.Password, user.Password)
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
	assert.Equal(t, updatedUser.UserID, findUser.UserID)
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

func TestRead(t *testing.T) {
	var i UserInteractor
	findUser, err := i.GetUserByEmail(DemoUser.Email)
	assert.Equal(t, nil, err)

	user, err := i.Read(findUser.UserID)
	assert.Equal(t, nil, err)
	assert.Equal(t, updatedName, user.UserName)
	assert.Equal(t, DemoUser.Email, user.Email)

}

func TestReadByIDNotExists(t *testing.T) {
	var i UserInteractor
	var searchID int32 = 10000

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
	assert.Equal(t, email, auth.Email)

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
	deleteUser.UserID = 10000

	err = i.Delete(deleteUser)
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

	err = i.Delete(deleteUser)
	assert.Equal(t, nil, err)

	findUser, err = i.GetUserByEmail(testemail)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, int32(0), findUser.UserID)
	assert.Equal(t, "", findUser.UserName)
	assert.Equal(t, "", findUser.Email)
}

func TestTokenAuth(t *testing.T) {
	var i UserInteractor
	initUserTable()
	user := &DemoUser
	user.Password = testpassword
	createUser, err := i.Create(user)
	assert.Equal(t, nil, err)

	auth, err := i.LoginAuth(createUser.Email, testpassword)
	assert.Equal(t, nil, err)

	authUser, err := i.TokenAuth(auth.Token)
	assert.Equal(t, nil, err)
	assert.Equal(t, createUser, authUser)
}

// List

// Search

func initUserTable() {
	DB := db.GetDB()
	DB.Delete(&model.User{})
}
