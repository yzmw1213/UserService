package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"github.com/go-playground/validator/v10"
	// grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/yzmw1213/UserService/domain/model"
	"github.com/yzmw1213/UserService/grpc/user_grpc"
	"github.com/yzmw1213/UserService/usecase/interactor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// StatusCreateUserSuccess ユーザ作成成功ステータス
	StatusCreateUserSuccess string = "USER_CREATE_SUCCESS"
	// StatusEmailAlreadyUsed 既に使われているEmail登録時のエラーステータス
	StatusEmailAlreadyUsed string = "EMAIL_ALREADY_USED_ERROR"
	// StatusEmailInputInvalid 無効なEmail入力時のエラーステータス
	StatusEmailInputInvalid string = "EMAIL_INPUT_INVALID_ERROR"
	// StatusUsernameNumError 無効な文字数Username入力時のエラーステータス
	StatusUsernameNumError string = "USERNAME_NUM_ERROR"
)

type server struct {
	Usecase interactor.UserInteractor
}

// ErrorCode エラーコード
var ErrorCode string

// NewUserGrpcServer gRPCサーバー起動
func NewUserGrpcServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	userServer := &server{}

	// s := grpc.NewServer(
	// 	grpc.UnaryInterceptor(
	// 		grpc_middleware.ChainUnaryServer(
	// 			grpc_auth.UnaryServerInterceptor(authorization.AuthFunc),
	// 		)),
	// )

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	user_grpc.RegisterUserServiceServer(s, userServer)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("main grpc server has started")

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a sgnal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the client")
	lis.Close()
	fmt.Println("End of Program")

}

// https://github.com/grpc/grpc-go/issues/106
func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	// ここでLoginAuthを呼び出す?
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// if len(md["username"]) > 0 && md["username"][0] == "admin" &&
		// 	len(md["password"]) > 0 && md["password"][0] == "admin123" {
		// 	return nil
		// }
		s := server{}
		_, err := s.Usecase.LoginAuth(md["email"][0], md["password"][0])
		if err != nil {
			return err
		}

		// Authをcontextに格納?
	}
	st := status.New(codes.Unauthenticated, "not authenticated")

	dt, err := st.WithDetails(
		&errdetails.LocalizedMessage{
			Locale:  "ja-JP",
			Message: "認証に失敗しました",
		},
		&errdetails.LocalizedMessage{
			Locale:  "en-US",
			Message: "Unauthenticaticated",
		},
	)
	if err != nil {
		return err
	}
	return dt.Err()
	// return EmptyMetadataErr
	// return fmt.Errorf("authorize Error has happen: Email: %v", mail)
}

func (s server) CreateUser(ctx context.Context, req *user_grpc.CreateUserRequest) (*user_grpc.CreateUserResponse, error) {
	postData := req.GetUser()
	ErrorCode = ""

	user := makeModel(postData)

	// 既に同一のemailによる登録がないかチェック
	if s.userExistsByEmail(user.Email) == true {
		return s.makeCreateResponse(nil, StatusEmailAlreadyUsed), nil
	}

	if _, err := s.Usecase.Create(user); err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			ErrorCode = makeValidationStatus(err)
			if ErrorCode != "" {
				break
			}
		}

		return s.makeCreateResponse(nil, ErrorCode), nil
	}
	// No errors
	return s.makeCreateResponse(postData, StatusCreateUserSuccess), nil
}

func (s server) DeleteUser(ctx context.Context, req *user_grpc.DeleteUserRequest) (*user_grpc.DeleteUserResponse, error) {
	postData := req.GetUserId()
	user := &model.User{
		UserID: postData,
	}
	if err := s.Usecase.Delete(user); err != nil {
		return nil, err
	}
	res := &user_grpc.DeleteUserResponse{}
	return res, nil
}

func (s server) ListUser(req *user_grpc.ListUserRequest, stream user_grpc.UserService_ListUserServer) error {
	rows, err := s.Usecase.List()
	if err != nil {
		return err
	}
	for _, user := range rows {
		user := &user_grpc.User{
			UserId:   user.UserID,
			UserName: user.UserName,
			Email:    user.Email,
		}
		res := &user_grpc.ListUserResponse{
			User: user,
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			log.Fatalf("Error while sending response to client :%v", sendErr)
			return sendErr
		}
	}

	return nil
}

func (s server) ReadUser(ctx context.Context, req *user_grpc.ReadUserRequest) (*user_grpc.ReadUserResponse, error) {
	userID := req.GetUserId()
	// userName := req.GetUserName()
	row, err := s.Usecase.Read(userID)
	if err != nil {
		return nil, err
	}
	user := &user_grpc.User{
		UserId: row.UserID,
	}
	res := &user_grpc.ReadUserResponse{
		User: user,
	}
	return res, nil
}

func (s server) UpdateUser(ctx context.Context, req *user_grpc.UpdateUserRequest) (*user_grpc.UpdateUserResponse, error) {
	postData := req.GetUser()
	validate := validator.New()

	user := makeModel(req.GetUser())

	// User構造体のバリデーション
	if error := validate.Struct(user); error != nil {
		return nil, error
	}

	if _, err := s.Usecase.Update(user); err != nil {
		return nil, err
	}
	res := &user_grpc.UpdateUserResponse{
		User: postData,
	}
	return res, nil
}

// userExistsByEmail Emailが登録されているユーザーが存在するかの判定
func (s server) userExistsByEmail(email string) bool {
	if email == "" {
		return false
	}
	user, _ := s.Usecase.GetUserByEmail(email)
	if user.UserID == 0 {
		return false
	}
	return true
}

func makeModel(gUser *user_grpc.User) *model.User {
	user := &model.User{
		UserID:   gUser.GetUserId(),
		UserName: gUser.GetUserName(),
		Password: gUser.GetPassword(),
		Email:    gUser.GetEmail(),
	}
	return user
}

func makeGrpcUser(user *model.User) *user_grpc.User {
	gUser := &user_grpc.User{
		UserId:   user.UserID,
		UserName: user.UserName,
		Password: user.Password,
		Email:    user.Password,
	}
	return gUser
}

func makeValidationStatus(err validator.FieldError) string {
	var code string = ""
	var field string = err.Field()
	var validationTag string = err.ActualTag()

	// Emailバリデーションエラー時
	if field == "Email" {
		code = StatusEmailInputInvalid
	} else if field == "UserName" {
		// UserNameの文字数が不適の場合
		for _, v := range []string{"min", "max"} {
			if v == validationTag {
				code = StatusUsernameNumError
				break
			}
		}
	} else {
		code = "unexpected error"
	}

	return code
}

// ログイン
// Email, Passwordの組み合わせで認証を行う
func (s server) Login(ctx context.Context, req *user_grpc.LoginRequest) (*user_grpc.LoginResponse, error) {
	if s.userExistsByEmail(req.GetEmail()) != true {
		return s.makeLoginResponse("", ""), nil
	}
	auth, err := s.Usecase.LoginAuth(req.GetEmail(), req.GetPassword())
	if err != nil {
		return s.makeLoginResponse("", ""), err
	}

	return s.makeLoginResponse(auth.Token, auth.Email), nil
}

// makeCreateResponse CreateUserメソッドのresponseを生成し返す
func (s server) makeCreateResponse(user *user_grpc.User, statusCode string) *user_grpc.CreateUserResponse {
	res := &user_grpc.CreateUserResponse{}
	if user != nil {
		res.User = user
	}
	if statusCode != "" {
		responseStatus := &user_grpc.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeLoginResponse CLoginメソッドのresponseを生成し返す
func (s server) makeLoginResponse(token string, email string) *user_grpc.LoginResponse {
	return &user_grpc.LoginResponse{
		Token: token,
		Email: email,
	}
}
