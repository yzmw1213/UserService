package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/UserService/db"
	"github.com/yzmw1213/UserService/domain/model"
	"github.com/yzmw1213/UserService/grpc/userservice"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize        = 1024 * 1024
	one     uint32 = 1
	nine    uint32 = 9
)

type loginCreds struct {
	Email, Password string
}

var lis *bufconn.Listener
var err error
var ctx = context.Background()

var demoUser = userservice.User{
	UserName:    "デモユーザ名1",
	Password:    "demopassword",
	ProfileText: "プロフィールが入ります",
	Authority:   one,
}

var demoSuperUser = userservice.User{
	UserName:    "manager",
	Email:       "super@gmail.com",
	Password:    "superpassword",
	ProfileText: "プロフィールが入ります",
	Authority:   nine,
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := makeServer()
	userservice.RegisterUserServiceServer(s, &server{})
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
	initUserTable()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)
	_, err = createDefaultUser(client, "demo@example.com", one)
	assert.Equal(t, nil, err)
	_, err = createDefaultUser(client, "wqrh3tws8@example.com", one)
	assert.Equal(t, nil, err)
	_, err = createDefaultUser(client, "rewthet3a@example.com", one)
	assert.Equal(t, nil, err)
	_, err = createDefaultUser(client, "super@example.com", nine)
	assert.Equal(t, nil, err)
}

// TestCreateUserEmailInvalid ユーザ作成Email無効の異常系
func TestCreateUserEmailInvalid(t *testing.T) {
	var createUser *userservice.User
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	createUser = &userservice.User{
		UserName: "テストユーザ名1",
		Password: "demopassword1",
		Email:    "aaaaa",
	}

	req := &userservice.CreateUserRequest{
		User: createUser,
	}

	_, err = client.CreateUser(ctx, req)
	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "Email", f)
	assert.Equal(t, StatusEmailInputInvalid, d)
}

// TestCreateUserEmailNull ユーザ作成Email未入力の異常系
func TestCreateUserEmailNull(t *testing.T) {
	var createUser *userservice.User
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	createUser = &userservice.User{
		UserName: "テストユーザ名2",
		Password: "demopassword",
		Email:    "",
	}

	req := &userservice.CreateUserRequest{
		User: createUser,
	}

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "Email", f)
	assert.Equal(t, StatusEmailInputInvalid, d)
}

// TestCreateUserEmailUsed 既に登録済みのemailを登録する異常系
func TestCreateUserEmailUsed(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	firstCreateUser := &userservice.User{
		UserName: "テストユーザ名3",
		Password: "demopassword3",
		Email:    "used@gmail.com",
	}

	nextCreateUser := &userservice.User{
		UserName: "テストユーザ名4",
		Password: "demopassword4",
		Email:    "used@gmail.com",
	}

	req := &userservice.CreateUserRequest{
		User: firstCreateUser,
	}

	_, err = client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)

	req = &userservice.CreateUserRequest{
		User: nextCreateUser,
	}

	// firstCreateUserで登録したemailを再び使用して登録
	res, err := client.CreateUser(ctx, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusEmailAlreadyUsed, res.GetStatus().GetCode())
}

// TestCreateUserNameNull ユーザ作成UserName未入力の異常系
func TestCreateUserNameNull(t *testing.T) {
	var createUser *userservice.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	createUser = &userservice.User{
		UserName: "",
		Password: "demopassword",
		Email:    "bbbbb@gmail.com",
	}

	req := &userservice.CreateUserRequest{
		User: createUser,
	}

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "UserName", f)
	assert.Equal(t, StatusUserNameCountError, d)
}

// TestCreateUserNameTooShort ユーザ作成UserName文字数不足の異常系
func TestCreateUserNameTooShort(t *testing.T) {
	var createUser *userservice.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	createUser = &userservice.User{
		UserName: "de",
		Password: "demopassword",
		Email:    "tooshort@gmail.com",
	}

	req := &userservice.CreateUserRequest{
		User: createUser,
	}

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "UserName", f)
	assert.Equal(t, StatusUserNameCountError, d)
}

// TestCreateUserNameTooLong ユーザ作成UserName文字数超過の異常系
func TestCreateUserNameTooLong(t *testing.T) {
	var createUser *userservice.User

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	createUser = &userservice.User{
		UserName: "demouserdemouserd",
		Password: "demopassword",
		Email:    "toolong@gmail.com",
	}

	req := &userservice.CreateUserRequest{
		User: createUser,
	}

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "UserName", f)
	assert.Equal(t, StatusUserNameCountError, d)
}

// TestLogin
func TestLogin(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	req := &userservice.LoginRequest{
		Email:    "demo@example.com",
		Password: "demopassword",
	}

	res, err := client.Login(ctx, req)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", res.GetAuth().GetToken())
	assert.NotEqual(t, zero, res.GetAuth().GetUserId())
	assert.NotEqual(t, "", res.GetUser().GetUserName())
	assert.NotEqual(t, zero, res.GetUser().GetUserId())
}

// TestLoginByWrongPassword 登録のないパスワードで認証を行う異常系
func TestLoginByNotRegistered(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	req := &userservice.LoginRequest{
		Email:    "demo@gmail.com",
		Password: "wrongpassword",
	}

	_, err = client.Login(ctx, req)
	assert.NotEqual(t, nil, err)
}

// TestGuestLogin ゲストログイン
func TestGuestLogin(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	req := &userservice.GuestLoginRequest{}

	res, err := client.GuestLogin(ctx, req)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", res.GetAuth().GetToken())
	assert.NotEqual(t, zero, res.GetAuth().GetUserId())
	assert.NotEqual(t, "", res.GetUser().GetUserName())
	assert.NotEqual(t, zero, res.GetUser().GetUserId())
}

// SuperUserLogin 管理ユーザーログイン
func SuperUserLogin(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	req := &userservice.SuperUserLoginRequest{}

	res, err := client.SuperUserLogin(ctx, req)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", res.GetAuth().GetToken())
	assert.NotEqual(t, zero, res.GetAuth().GetUserId())
	assert.NotEqual(t, "", res.GetUser().GetUserName())
	assert.NotEqual(t, zero, res.GetUser().GetUserId())
}

func getErrorDetail(err error) (string, string) {
	var field string
	var description string
	st, _ := status.FromError(err)
	for _, detail := range st.Details() {
		switch dType := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range dType.GetFieldViolations() {
				field = violation.GetField()
				description = violation.GetDescription()
			}
		}
	}

	return field, description
}

// TestAuth トークンよりユーザー情報を返す
func TestAuth(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	loginReq := &userservice.LoginRequest{
		Email:    "demo@example.com",
		Password: demoUser.Password,
	}

	loginRes, err := client.Login(ctx, loginReq)
	token := loginRes.GetAuth().GetToken()
	userID := loginRes.GetAuth().GetUserId()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", token)
	assert.NotEqual(t, zero, userID)

	tokenAuthReq := &userservice.TokenAuthRequest{
		Token: token,
	}

	tokenAuthRes, err := client.TokenAuth(ctx, tokenAuthReq)
	assert.Equal(t, nil, err)
	assert.Equal(t, "demo@example.com", tokenAuthRes.GetUser().GetEmail())
}

func TestFollowUser(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	followReq := &userservice.FollowUserRequet{
		FollwerUserId: 1,
		FollwedUserId: 2,
	}

	res, err := client.FollowUser(ctx, followReq)
	assert.Equal(t, StatusFollowSuccess, res.GetStatus().GetCode())
	assert.Equal(t, nil, err)
}

func TestUnFollowUser(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userservice.NewUserServiceClient(conn)

	unfollowReq := &userservice.UnFollowUserRequet{
		FollwerUserId: 1,
		FollwedUserId: 2,
	}

	res, err := client.UnFollowUser(ctx, unfollowReq)
	assert.Equal(t, StatusUnFollowSuccess, res.GetStatus().GetCode())
	assert.Equal(t, nil, err)
}

func createDefaultUser(c userservice.UserServiceClient, email string, authority uint32) (*userservice.CreateUserResponse, error) {
	user := &demoUser
	user.Email = email
	user.Authority = authority

	req := &userservice.CreateUserRequest{
		User: user,
	}

	return c.CreateUser(ctx, req)
}

func initUserTable() {
	DB := db.GetDB()
	DB.Delete(&model.User{})
	DB.Delete(&model.Relation{})
}
