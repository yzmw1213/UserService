package authorization

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pkg/errors"
	"github.com/yzmw1213/UserService/domain/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultAuthenticateFunc はすべてのサービスに共通する認証処理を行う関数を表す。
type DefaultAuthenticateFunc func(ctx context.Context) (context.Context, error)

type key int

const (
	expectedScheme = "bearer"
	// TokenKey トークンキー
	TokenKey = "token"
	// secretを環境変数から読む
	secret = "2FMd5FNSqS/nW2wWJy5S3ppjSHhUnLt8HuwBkTD6HqfPfBBDlykwLA=="
	//
	strigKey key = iota
	// ゼロ値
	zero uint32 = 0
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

	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"could not read auth token: %v",
			err,
		)
	}
	//  Emailをtokenに格納
	claims := &jwt.MapClaims{
		// "email": ctx.Value("email"),
		"id": ctx.Value("id"),
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
	return context.WithValue(ctx, strigKey, token)
}

// CreateToken ユーザー情報からトークンを発行する
func CreateToken(user *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"id": user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

// ParseToken は jwt トークンから元になった認証情報を取り出す。
func ParseToken(signedString string) (uint32, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return zero, errors.Wrapf(err, "%s is expired", signedString)
			}

			return zero, errors.Wrapf(err, "%s is invalid", signedString)
		}

		return zero, errors.Wrapf(err, "%s is invalid", signedString)
	}
	if token == nil {
		return zero, fmt.Errorf("not found token in %s", signedString)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return zero, fmt.Errorf("not found claims in %s", signedString)
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return zero, fmt.Errorf("not found claims in %s", signedString)
	}

	return uint32(userID), nil
}
