package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yzmw1213/UserService/grpc/user_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize = 1024 * 1024
)

type loginCreds struct {
	Email, Password string
}

var lis *bufconn.Listener
var err error
var ctx = context.Background()

var demouser = user_grpc.User{
	UserName: "デモユーザ名1",
	Email:    "demo@gmail.com",
	Password: "demopassword",
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	user_grpc.RegisterUserServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Email,
		"password": c.Password,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}

// func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
//     return map[string]string{
//         "username": c.Username,
//         "password": c.Password,
//     }, nil
// }

// TestCreateUser ユーザ作成正常系
func TestCreateUser(t *testing.T) {
	var createUsers []*user_grpc.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUsers = append(createUsers, &demouser)

	for _, user := range createUsers {
		req := &user_grpc.CreateUserRequest{
			User: user,
		}

		res, err := client.CreateUser(ctx, req)

		assert.Equal(t, nil, err)
		assert.Equal(t, StatusCreateUserSuccess, res.GetStatus().GetCode())
	}
}

// TestCreateUserEmailInvalidWithAuth ユーザ作成Email無効の異常系
func TestCreateUserEmailInvalidWithAuth(t *testing.T) {
	var createUser *user_grpc.User
	conn, err := grpc.Dial(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&loginCreds{
			Email:    "demo@gmail.com",
			Password: "demopassword",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "テストユーザ名1",
		Password: "demopassword1",
		Email:    "aaaaa",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusEmailInputInvalid, res.GetStatus().GetCode())
}

// TestCreateUserEmailInvalid ユーザ作成Email無効の異常系
func TestCreateUserEmailInvalid(t *testing.T) {
	var createUser *user_grpc.User
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "テストユーザ名1",
		Password: "demopassword1",
		Email:    "aaaaa",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusEmailInputInvalid, res.GetStatus().GetCode())
}

// TestCreateUserEmailNull ユーザ作成Email未入力の異常系
func TestCreateUserEmailNull(t *testing.T) {
	var createUser *user_grpc.User
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "テストユーザ名2",
		Password: "demopassword",
		Email:    "",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusEmailInputInvalid, res.GetStatus().GetCode())
}

// TestCreateUserEmailUsed 既に登録済みのemailを登録する異常系
func TestCreateUserEmailUsed(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	firstCreateUser := &user_grpc.User{
		UserName: "テストユーザ名3",
		Password: "demopassword3",
		Email:    "used@gmail.com",
	}

	nextCreateUser := &user_grpc.User{
		UserName: "テストユーザ名4",
		Password: "demopassword4",
		Email:    "used@gmail.com",
	}

	req := &user_grpc.CreateUserRequest{
		User: firstCreateUser,
	}

	_, err = client.CreateUser(ctx, req)

	if err != nil {
		log.Fatalf("Error message for Email")
	}

	req = &user_grpc.CreateUserRequest{
		User: nextCreateUser,
	}

	// firstCreateUserで登録したemailを再び使用して登録
	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusEmailAlreadyUsed, res.GetStatus().GetCode())
}

// TestCreateUserNameNull ユーザ作成UserName未入力の異常系
func TestCreateUserNameNull(t *testing.T) {
	var createUser *user_grpc.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "",
		Password: "demopassword",
		Email:    "bbbbb@gmail.com",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusUsernameNumError, res.GetStatus().GetCode())
}

// TestCreateUserNameTooShort ユーザ作成UserName文字数不足の異常系
func TestCreateUserNameTooShort(t *testing.T) {
	var createUser *user_grpc.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "demou",
		Password: "demopassword",
		Email:    "tooshort@gmail.com",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusUsernameNumError, res.GetStatus().GetCode())
}

// TestCreateUserNameTooLong ユーザ作成UserName文字数超過の異常系
func TestCreateUserNameTooLong(t *testing.T) {
	var createUser *user_grpc.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	createUser = &user_grpc.User{
		UserName: "demouserdemouserd",
		Password: "demopassword",
		Email:    "toolong@gmail.com",
	}

	req := &user_grpc.CreateUserRequest{
		User: createUser,
	}

	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusUsernameNumError, res.GetStatus().GetCode())
}

// TestLogin
func TestLogin(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	req := &user_grpc.LoginRequest{
		Email:    "demo@gmail.com",
		Password: "demopassword",
	}

	res, err := client.Login(ctx, req)
	assert.Equal(t, nil, err)
	assert.Equal(t, "demo@gmail.com", res.GetEmail())
	assert.NotEqual(t, "", res.GetToken())
}

// TestTokenAuth
// func TestTokenAuth(t *testing.T) {
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer conn.Close()

// 	client := user_grpc.NewUserServiceClient(conn)

// 	loginReq := &user_grpc.LoginRequest{
// 		Email:    "demo@gmail.com",
// 		Password: "demopassword",
// 	}

// 	loginRes, err := client.Login(ctx, loginReq)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, "demo@gmail.com", loginRes.GetEmail())
// 	assert.NotEqual(t, "", loginRes.GetToken())

// 	tokenAuthReq :=

// }
