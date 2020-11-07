package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/UserService/grpc/user_grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize        = 1024 * 1024
	one     uint32 = 1
	two     uint32 = 2
	nine    uint32 = 9
)

type loginCreds struct {
	Email, Password string
}

var lis *bufconn.Listener
var err error
var ctx = context.Background()

var demoUser = user_grpc.User{
	UserName:    "デモユーザ名1",
	Email:       "demo@gmail.com",
	Password:    "demopassword",
	ProfileText: "プロフィールが入ります",
	Authority:   one,
}

var demoSuperUser = user_grpc.User{
	UserName:    "manager",
	Email:       "super@gmail.com",
	Password:    "superpassword",
	ProfileText: "プロフィールが入ります",
	Authority:   nine,
}

var demoCompanyUser = user_grpc.User{
	UserName:    "companyA",
	Email:       "company@gmail.com",
	Password:    "companypassword",
	ProfileText: "プロフィールが入ります",
	Authority:   two,
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := makeServer()
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

	createUsers = append(createUsers, &demoUser)
	createUsers = append(createUsers, &demoSuperUser)
	createUsers = append(createUsers, &demoCompanyUser)

	for _, user := range createUsers {
		req := &user_grpc.CreateUserRequest{
			User: user,
		}

		res, err := client.CreateUser(ctx, req)

		assert.Equal(t, nil, err)
		assert.Equal(t, StatusCreateUserSuccess, res.GetStatus().GetCode())
	}
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

	_, err = client.CreateUser(ctx, req)
	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "Email", f)
	assert.Equal(t, StatusEmailInputInvalid, d)
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

	assert.Equal(t, nil, err)

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

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "UserName", f)
	assert.Equal(t, StatusUserNameCountError, d)
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

	_, err = client.CreateUser(ctx, req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "UserName", f)
	assert.Equal(t, StatusUserNameCountError, d)
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

	client := user_grpc.NewUserServiceClient(conn)

	req := &user_grpc.LoginRequest{
		Email:    "demo@gmail.com",
		Password: "demopassword",
	}

	res, err := client.Login(ctx, req)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", res.GetAuth().GetToken())
	assert.NotEqual(t, zero, res.GetAuth().GetUserId())
}

// TestLoginByWrongPassword 登録のないパスワードで認証を行う異常系
func TestLoginByNotRegistered(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := user_grpc.NewUserServiceClient(conn)

	req := &user_grpc.LoginRequest{
		Email:    "demo@gmail.com",
		Password: "wrongpassword",
	}

	_, err = client.Login(ctx, req)
	assert.NotEqual(t, nil, err)
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

	client := user_grpc.NewUserServiceClient(conn)

	loginReq := &user_grpc.LoginRequest{
		Email:    demoUser.Email,
		Password: demoUser.Password,
	}

	loginRes, err := client.Login(ctx, loginReq)
	token := loginRes.GetAuth().GetToken()
	userID := loginRes.GetAuth().GetUserId()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", token)
	assert.NotEqual(t, zero, userID)

	tokenAuthReq := &user_grpc.TokenAuthRequest{
		Token: token,
	}

	tokenAuthRes, err := client.TokenAuth(ctx, tokenAuthReq)
	assert.Equal(t, nil, err)
	assert.Equal(t, demoUser.Email, tokenAuthRes.GetUser().GetEmail())

}
