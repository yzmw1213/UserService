package authorization

import (
	"context"
	"log"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultAuthenticateFunc はすべてのサービスに共通する認証処理を行う関数を表す。
type DefaultAuthenticateFunc func(ctx context.Context) (context.Context, error)

const (
	expectedScheme = "bearer"
	// TokenKey トークンキー
	TokenKey = "token"
)

// UnaryServerInterceptor はリクエストごとの認証処理を行う、unary サーバーインターセプターを返す。
func UnaryServerInterceptor(authFunc DefaultAuthenticateFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := authFunc(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return handler(newCtx, req)
	}
}

// AuthFunc 署名付きトークンを生成する
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, expectedScheme)
	log.Println("AuthFunc")

	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"could not read auth token: %v",
			err,
		)

	}

	//  Emailをtokenに格納
	claims := &jwt.MapClaims{
		"email": ctx.Value("email"),
	}
	parser := new(jwt.Parser)
	parsedToken, _, err := parser.ParseUnverified(token, claims)

	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"could not parsed auth token: %v",
			err,
		)
	}

	return setToken(ctx, parsedToken.Claims.(*jwt.MapClaims)), nil

}

func setToken(ctx context.Context, token *jwt.MapClaims) context.Context {
	return context.WithValue(ctx, TokenKey, token)
}

// func GetToken(ctx context.Context) *jwt.MapClaims {
// 	return ctx.Value(TokenKey).(*jwt.MapClaims)
// }
