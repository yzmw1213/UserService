package interactor

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/yzmw1213/UserService/db"
	"github.com/yzmw1213/UserService/domain/model"
	"github.com/yzmw1213/UserService/usecase/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	err      error
	user     model.User
	users    []model.User
	rows     *sql.Rows
	validate *validator.Validate
)

// secretを環境変数から読む
const secret = "2FMd5FNSqS/nW2wWJy5S3ppjSHhUnLt8HuwBkTD6HqfPfBBDlykwLA=="

// UserInteractor ユーザサービスを提供するメソッド群
type UserInteractor struct{}

var _ repository.UserRepository = (*UserInteractor)(nil)

// Create ユーザ1件を作成
func (i *UserInteractor) Create(postData *model.User) (*model.User, error) {
	validate = validator.New()
	DB := db.GetDB()
	createUser := postData

	// User構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}
	inputPassword := postData.Password

	hash, err := createHashPassword(inputPassword)
	createUser.Password = hash

	if err != nil {
		return createUser, err
	}

	if err := DB.Create(createUser).Error; err != nil {
		return postData, err
	}

	return postData, nil
}

// Delete ユーザ1件を削除
func (i *UserInteractor) Delete(postData *model.User) error {
	DB := db.GetDB()
	if err := DB.Delete(postData).Error; err != nil {
		return err
	}
	return nil
}

// Count ユーザ件数を取得
func (i *UserInteractor) Count(user model.User) (int, error) {
	var count int
	DB := db.GetDB()
	if err := DB.Find(&user).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// List ユーザを全件取得
func (i *UserInteractor) List() ([]model.User, error) {
	var userList []model.User
	rows, err := listAll(context.Background())
	if err != nil {
		fmt.Println("Error happened")
		return []model.User{}, err
	}
	for _, row := range rows {
		userList = append(userList, row)
	}

	return userList, nil
}

// listAll 全件取得
func listAll(ctx context.Context) ([]model.User, error) {
	DB := db.GetDB()

	rows, err := DB.Find(&users).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		DB.ScanRows(rows, &user)
		users = append(users, user)
	}
	return users, nil
}

// Update ユーザを更新する
func (i *UserInteractor) Update(postData *model.User) (*model.User, error) {
	DB := db.GetDB()
	// postされたIdに紐づくuserを取得
	id := postData.UserID
	findUser := &model.User{}

	// User構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}

	if err := DB.Where("user_id = ?", id).First(&findUser).Error; err != nil {
		log.Fatalf("err: %v", err)
		return findUser, err
	}

	updateUser := postData
	// パスワードをhash
	hash, err := createHashPassword(postData.Password)
	// hashしたパスワードをSaveするuserにセット
	updateUser.Password = string(hash)

	if err != nil {
		return updateUser, err
	}

	updateUser.UserID = findUser.UserID

	if err := DB.Save(updateUser).Error; err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

// Read IDを元にユーザを1件取得する
func (i *UserInteractor) Read(userID int32) (model.User, error) {
	DB := db.GetDB()
	row := DB.First(&user, userID)
	if err := row.Error; err != nil {
		return model.User{}, err
	}
	DB.Table(db.TableName).Scan(row)
	return user, nil
}

// GetUserByEmail Emailを元にユーザを1件取得する
func (i *UserInteractor) GetUserByEmail(email string) (model.User, error) {
	var user model.User

	DB := db.GetDB()
	row := DB.Where("email = ?", email).First(&user)
	if err := row.Error; err != nil {
		return user, err
	}
	DB.Table(db.TableName).Scan(row)

	return user, nil
}

func createHashPassword(password string) (string, error) {
	// パスワードの暗号化
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashPassword := string(hash)

	if err != nil {
		log.Fatal(err)
		return hashPassword, err
	}
	return hashPassword, nil
}

// LoginAuth パスワード入力による認証メソッド
func (i *UserInteractor) LoginAuth(email string, inputPassword string) (*model.Auth, error) {
	// 入力ユーザ存在有無の判定.
	// eメールに紐づくユーザのパスワードを取得
	findUser, err := i.GetUserByEmail(email)
	if err != nil {
		return &model.Auth{}, err
	}

	// DBから取得したパスワードと入力値のハッシュを比較
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(inputPassword))
	// 認証失敗
	if err != nil {
		return &model.Auth{}, err
	}

	// contextにユーザ情報格納
	ctx := context.Background()
	ctx = context.WithValue(ctx, "email", findUser.Email)

	// authのAuthFuncを呼び出す
	// jwt生成
	token, err := createToken(&findUser)
	if err != nil {
		return &model.Auth{}, err
	}

	return &model.Auth{
		Token: token,
		Email: email,
	}, nil
}

// TokenAuth 認証トークンで認証を行い、ユーザ情報を返す
func (i *UserInteractor) TokenAuth(token string) (model.User, error) {
	// tokenからemailを取り出す
	email, err := parseToken(token)
	if err != nil {
		return model.User{}, err
	}
	// emailでユーザ検索
	user, err := i.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	//
	return user, nil
}

// parseToken は jwt トークンから元になった認証情報を取り出す。
func parseToken(signedString string) (string, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	log.Println(signedString)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.Wrapf(err, "%s is expired", signedString)

			}

			return "", errors.Wrapf(err, "%s is invalid", signedString)
		}

		return "", errors.Wrapf(err, "%s is invalid", signedString)
	}

	if token == nil {

		return "", fmt.Errorf("not found token in %s", signedString)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("not found claims in %s", signedString)
	}
	email := claims["email"].(string)

	return email, nil
}

func createToken(user *model.User) (string, error) {

	claims := &jwt.MapClaims{
		"email": user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}
