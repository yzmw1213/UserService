package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/yzmw1213/UserService/authorization"

	"github.com/yzmw1213/UserService/domain/model"
	"github.com/yzmw1213/UserService/grpc/user_grpc"
)

const (
	// StatusCreateUserSuccess ユーザ作成成功ステータス
	StatusCreateUserSuccess string = "USER_CREATE_SUCCESS"
	// StatusUpdateUserSuccess ユーザ更新成功ステータス
	StatusUpdateUserSuccess string = "USER_UPDATE_SUCCESS"
	// StatusEmailAlreadyUsed 既に使われているEmail登録時のエラーステータス
	StatusEmailAlreadyUsed string = "EMAIL_ALREADY_USED_ERROR"
	// StatusEmailInputInvalid 無効なEmail入力時のエラーステータス
	StatusEmailInputInvalid string = "EMAIL_INPUT_INVALID_ERROR"
	// StatusUserNameCountError 無効な文字数Username入力時のエラーステータス
	StatusUserNameCountError string = "USERNAME_NUM_ERROR"
	// zero ユーザーIDのゼロ値
	zero uint32 = 0
)

func (s server) CreateUser(ctx context.Context, req *user_grpc.CreateUserRequest) (*user_grpc.CreateUserResponse, error) {
	user := makeModel(req.GetUser())

	// 既に同一のemailによる登録がないかチェック
	if s.userExistsByEmail(user.Email) == true {
		return s.makeCreateUserResponse(StatusEmailAlreadyUsed), nil
	}

	// clientのリクエストに使うため、token作成をhandlerに出す?
	user, err := s.Usecase.Create(user)

	if err != nil {
		return nil, err
	}
	// ユーザー作成に成功したら、clientより他のサービスのメソッドにリクエストを送る、など

	// No errors
	return s.makeCreateUserResponse(StatusCreateUserSuccess), nil
}

func (s server) DeleteUser(ctx context.Context, req *user_grpc.DeleteUserRequest) (*user_grpc.DeleteUserResponse, error) {
	postData := req.GetUserId()
	user := &model.User{
		ID: postData,
	}
	if err := s.Usecase.DeleteByID(user.ID); err != nil {
		return nil, err
	}
	res := &user_grpc.DeleteUserResponse{}
	return res, nil
}

func (s server) ListUser(ctx context.Context, req *user_grpc.ListUserRequest) (*user_grpc.ListUserResponse, error) {
	rows, err := s.Usecase.List()
	if err != nil {
		return nil, err
	}
	var users []*user_grpc.User
	for _, user := range rows {
		user := makeGrpcUser(&user)
		users = append(users, user)
	}
	res := &user_grpc.ListUserResponse{
		User: users,
	}

	return res, nil
}

func (s server) ReadUser(ctx context.Context, req *user_grpc.ReadUserRequest) (*user_grpc.ReadUserResponse, error) {
	userID := req.GetUserId()
	// userName := req.GetUserName()
	row, err := s.Usecase.Read(userID)
	if err != nil {
		return nil, err
	}
	user := &user_grpc.User{
		UserId: row.ID,
	}
	res := &user_grpc.ReadUserResponse{
		User: user,
	}
	return res, nil
}

func (s server) UpdateUser(ctx context.Context, req *user_grpc.UpdateUserRequest) (*user_grpc.UpdateUserResponse, error) {
	user := makeModel(req.GetUser())

	if _, err := s.Usecase.Update(user); err != nil {
		return nil, err
	}

	return s.makeUpdateUserResponse(StatusUpdateUserSuccess), nil
}

// userExistsByEmail Emailが登録されているユーザーが存在するかの判定
func (s server) userExistsByEmail(email string) bool {
	if email == "" {
		return false
	}
	user, _ := s.Usecase.GetUserByEmail(email)
	if user.ID == 0 {
		return false
	}
	return true
}

func makeModel(gUser *user_grpc.User) *model.User {
	user := &model.User{
		ID:        gUser.GetUserId(),
		UserName:  gUser.GetUserName(),
		Password:  gUser.GetPassword(),
		Email:     gUser.GetEmail(),
		Authority: gUser.GetAuthority(),
	}

	return user
}

func makeGrpcUser(user *model.User) *user_grpc.User {
	gUser := &user_grpc.User{
		UserId:    user.ID,
		UserName:  user.UserName,
		Password:  user.Password,
		Email:     user.Email,
		Authority: user.Authority,
	}
	return gUser
}

func makeGrpcAuth(auth *model.Auth) *user_grpc.Auth {
	gAuth := &user_grpc.Auth{
		Token:     auth.Token,
		UserId:    auth.UserID,
		Authority: auth.Authority,
	}
	return gAuth
}

// ログイン
// Email, Passwordの組み合わせで認証を行う
func (s server) Login(ctx context.Context, req *user_grpc.LoginRequest) (*user_grpc.LoginResponse, error) {
	log.Println(req.Email)
	log.Println(req.Password)
	if s.userExistsByEmail(req.GetEmail()) != true {
		return s.makeLoginResponse(&user_grpc.Auth{Token: "", UserId: zero}), errors.New("user not found")
	}
	auth, err := s.Usecase.LoginAuth(req.GetEmail(), req.GetPassword())
	if err != nil {
		return s.makeLoginResponse(&user_grpc.Auth{Token: "", UserId: zero}), err
	}
	return s.makeLoginResponse(makeGrpcAuth(auth)), nil
}

func (s server) TokenAuth(ctx context.Context, req *user_grpc.TokenAuthRequest) (*user_grpc.TokenAuthResponse, error) {
	// tokenからemailを取り出す
	id, err := authorization.ParseToken(req.GetToken())
	if err != nil {
		return nil, err
	}
	// user, err := s.Usecase.GetUserByEmail(email)
	user, err := s.Usecase.GetUserByUserID(id)
	if err != nil {
		return nil, err
	}

	return &user_grpc.TokenAuthResponse{
		User: makeGrpcUser(&user),
	}, nil
}

// makeCreateUserResponse CreateUserメソッドのresponseを生成し返す
func (s server) makeCreateUserResponse(statusCode string) *user_grpc.CreateUserResponse {
	res := &user_grpc.CreateUserResponse{}
	if statusCode != "" {
		responseStatus := &user_grpc.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdateUserResponse UpdateUserメソッドのresponseを生成し返す
func (s server) makeUpdateUserResponse(statusCode string) *user_grpc.UpdateUserResponse {
	res := &user_grpc.UpdateUserResponse{}
	if statusCode != "" {
		responseStatus := &user_grpc.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeLoginResponse CLoginメソッドのresponseを生成し返す
func (s server) makeLoginResponse(auth *user_grpc.Auth) *user_grpc.LoginResponse {
	return &user_grpc.LoginResponse{
		Auth: auth,
	}
}
