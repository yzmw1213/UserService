package grpc

import (
	"context"
	"errors"

	"github.com/yzmw1213/UserService/authorization"

	"github.com/yzmw1213/UserService/domain/model"
	"github.com/yzmw1213/UserService/grpc/userservice"
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
	// StatusFollowSuccess フォロー成功ステータス
	StatusFollowSuccess string = "FOLLOW_SUCCESS"
	// StatusUnFollowSuccess フォロー解除成功ステータス
	StatusUnFollowSuccess string = "UNFOLLOW_SUCCESS"
	// StatusDeleteUserSuccess ユーザ削除成功ステータス
	StatusDeleteUserSuccess string = "USER_DELETE_SUCCESS"
	// zero ユーザーIDのゼロ値
	zero uint32 = 0
)

func (s server) CreateUser(ctx context.Context, req *userservice.CreateUserRequest) (*userservice.CreateUserResponse, error) {
	user := makeModel(req.GetUser())

	// 既に同一のemailによる登録がないかチェック
	if s.userExistsByEmail(zero, user.Email) == true {
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

func (s server) DeleteUser(ctx context.Context, req *userservice.DeleteUserRequest) (*userservice.DeleteUserResponse, error) {
	postData := req.GetUserId()
	user := &model.User{
		ID: postData,
	}
	if err := s.Usecase.DeleteByID(user.ID); err != nil {
		return nil, err
	}
	res := &userservice.DeleteUserResponse{
		Status: &userservice.ResponseStatus{
			Code: StatusDeleteUserSuccess,
		},
	}
	return res, nil
}

func (s server) ListUser(ctx context.Context, req *userservice.ListUserRequest) (*userservice.ListUserResponse, error) {
	rows, err := s.Usecase.ListAllNormalUser()
	if err != nil {
		return nil, err
	}
	var profiles []*userservice.UserProfile
	for _, user := range rows {
		user := makeGrpcUserProfile(&user, []uint32{})
		profiles = append(profiles, user)
	}
	res := &userservice.ListUserResponse{
		Profile: profiles,
	}

	return res, nil
}

func (s server) ReadUser(ctx context.Context, req *userservice.ReadUserRequest) (*userservice.ReadUserResponse, error) {
	userID := req.GetUserId()
	row, err := s.Usecase.GetUserByUserID(userID)
	followUsers := s.Usecase.GetFollowUsersByID(userID)

	if err != nil {
		return nil, err
	}
	user := makeGrpcUserProfile(&row, followUsers)
	res := &userservice.ReadUserResponse{
		Profile: user,
	}
	return res, nil
}

func (s server) UpdateUser(ctx context.Context, req *userservice.UpdateUserRequest) (*userservice.UpdateUserResponse, error) {
	user := makeModel(req.GetUser())

	// 既に同一のemailによる登録がないかチェック
	if s.userExistsByEmail(user.ID, user.Email) == true {
		return s.makeUpdateUserResponse(StatusEmailAlreadyUsed), nil
	}

	if _, err := s.Usecase.Update(user); err != nil {
		return nil, err
	}

	return s.makeUpdateUserResponse(StatusUpdateUserSuccess), nil
}

// userExistsByEmail Emailが登録されているユーザーが存在するかの判定
func (s server) userExistsByEmail(id uint32, email string) bool {
	if email == "" {
		return false
	}
	user, _ := s.Usecase.GetUserByEmail(email)
	if user.ID == 0 {
		return false
	}
	return user.ID > 0 && user.ID != id
}

func makeModel(gUser *userservice.User) *model.User {
	user := &model.User{
		ID:        gUser.GetUserId(),
		UserName:  gUser.GetUserName(),
		Password:  gUser.GetPassword(),
		Email:     gUser.GetEmail(),
		Authority: gUser.GetAuthority(),
	}

	return user
}

func makeGrpcUser(user *model.User) *userservice.User {
	gUser := &userservice.User{
		UserId:    user.ID,
		UserName:  user.UserName,
		Password:  user.Password,
		Email:     user.Email,
		Authority: user.Authority,
	}
	return gUser
}

func makeGrpcUserProfile(user *model.User, followUsers []uint32) *userservice.UserProfile {
	gUser := &userservice.UserProfile{
		UserId:      user.ID,
		UserName:    user.UserName,
		ProfileText: user.ProfileText,
		Authority:   user.Authority,
		FollowUsers: followUsers,
	}
	return gUser
}

func makeGrpcAuth(auth *model.Auth) *userservice.Auth {
	gAuth := &userservice.Auth{
		Token:     auth.Token,
		UserId:    auth.UserID,
		Authority: auth.Authority,
	}
	return gAuth
}

// ログイン
// Email, Passwordの組み合わせで認証を行う
func (s server) Login(ctx context.Context, req *userservice.LoginRequest) (*userservice.LoginResponse, error) {
	if s.userExistsByEmail(zero, req.GetEmail()) != true {
		return s.makeLoginResponse(&userservice.Auth{Token: "", UserId: zero}, &userservice.User{UserId: zero, UserName: ""}), errors.New("user not found")
	}
	auth, err := s.Usecase.LoginAuth(req.GetEmail(), req.GetPassword())
	user, err := s.Usecase.GetUserByUserID(auth.UserID)
	if err != nil {
		return s.makeLoginResponse(&userservice.Auth{Token: "", UserId: zero}, &userservice.User{UserId: zero, UserName: ""}), err
	}

	return s.makeLoginResponse(makeGrpcAuth(auth), makeGrpcUser(&user)), nil
}

// ゲストユーザーログイン
func (s server) GuestLogin(ctx context.Context, req *userservice.GuestLoginRequest) (*userservice.LoginResponse, error) {
	// ゲストアカウントを作成し、認証情報を取得
	auth, err := s.Usecase.CreateDemoUser()
	// 認証情報を元にユーザー情報を取得
	user, err := s.Usecase.GetUserByUserID(auth.UserID)

	if err != nil {
		return s.makeLoginResponse(&userservice.Auth{Token: "", UserId: zero}, &userservice.User{UserId: zero, UserName: ""}), err
	}

	return s.makeLoginResponse(makeGrpcAuth(auth), makeGrpcUser(&user)), nil
}

// 管理者ユーザーログイン
func (s server) SuperUserLogin(ctx context.Context, req *userservice.SuperUserLoginRequest) (*userservice.LoginResponse, error) {
	// 管理アカウントを作成し、認証情報を取得
	auth, err := s.Usecase.CreateDemoSuperUser()
	// 認証情報を元にユーザー情報を取得
	user, err := s.Usecase.GetUserByUserID(auth.UserID)

	if err != nil {
		return s.makeLoginResponse(&userservice.Auth{Token: "", UserId: zero}, &userservice.User{UserId: zero, UserName: ""}), err
	}

	return s.makeLoginResponse(makeGrpcAuth(auth), makeGrpcUser(&user)), nil
}

func (s server) TokenAuth(ctx context.Context, req *userservice.TokenAuthRequest) (*userservice.TokenAuthResponse, error) {
	// tokenからidを取り出す
	id, err := authorization.ParseToken(req.GetToken())
	if err != nil {
		return nil, err
	}
	user, err := s.Usecase.GetUserByUserID(id)
	if err != nil {
		return nil, err
	}

	return &userservice.TokenAuthResponse{
		User: makeGrpcUser(&user),
	}, nil
}

func (s server) FollowUser(ctx context.Context, req *userservice.FollowUserRequet) (*userservice.FollowUserResponse, error) {
	relation := &model.Relation{
		FollowerUserID: req.GetFollwerUserId(),
		FollowedUserID: req.GetFollwedUserId(),
	}

	if _, err := s.Usecase.Follow(relation); err != nil {
		return nil, err
	}

	return &userservice.FollowUserResponse{Status: &userservice.ResponseStatus{Code: StatusFollowSuccess}}, nil
}

func (s server) UnFollowUser(ctx context.Context, req *userservice.UnFollowUserRequet) (*userservice.UnFollowUserResponse, error) {
	relation := &model.Relation{
		FollowerUserID: req.GetFollwerUserId(),
		FollowedUserID: req.GetFollwedUserId(),
	}

	if _, err := s.Usecase.UnFollow(relation); err != nil {
		return nil, err
	}

	return &userservice.UnFollowUserResponse{Status: &userservice.ResponseStatus{Code: StatusUnFollowSuccess}}, nil
}

// makeCreateUserResponse CreateUserメソッドのresponseを生成し返す
func (s server) makeCreateUserResponse(statusCode string) *userservice.CreateUserResponse {
	res := &userservice.CreateUserResponse{}
	if statusCode != "" {
		responseStatus := &userservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdateUserResponse UpdateUserメソッドのresponseを生成し返す
func (s server) makeUpdateUserResponse(statusCode string) *userservice.UpdateUserResponse {
	res := &userservice.UpdateUserResponse{}
	if statusCode != "" {
		responseStatus := &userservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeLoginResponse Loginメソッドのresponseを生成し返す
func (s server) makeLoginResponse(auth *userservice.Auth, user *userservice.User) *userservice.LoginResponse {
	return &userservice.LoginResponse{
		Auth: auth,
		User: user,
	}
}

// makeErrorLoginResponse Loginメソッドがエラー時のresponseを生成し返す
func (s server) makeErrorLoginResponse() *userservice.LoginResponse {
	return s.makeLoginResponse(&userservice.Auth{Token: "", UserId: zero}, &userservice.User{UserId: zero, UserName: ""})
}
