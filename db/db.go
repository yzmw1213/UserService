package db

import (
	"fmt"
	"os"

	// gormのmysql接続用
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/UserService/domain/model"

	"github.com/jinzhu/gorm"
)

var (
	// DB データベース構造体
	DB *gorm.DB
	// tx トランザクション
	tx   *gorm.DB
	user model.User
	// UserTableName ユーザーサービステーブル名
	UserTableName string = "users"
	// RelationTableName フォロー関係テーブル名
	RelationTableName string = "relations"
)

func initDB() {
	var err error
	DBMS := "mysql"
	DBNAME := os.Getenv("DB_NAME")
	PASSWORD := os.Getenv("DB_PASSWORD")
	USER := os.Getenv("DB_USER")
	PROTOCOL := fmt.Sprintf("tcp(%s)", os.Getenv("DB_ADRESS"))
	OPTION := "?charset=utf8mb4&parseTime=True&loc=Local"
	CONNECTION := fmt.Sprintf("%s:%s@%s/%s%s", USER, PASSWORD, PROTOCOL, DBNAME, OPTION)

	DB, err = gorm.Open(DBMS, CONNECTION)
	if err != nil {
		panic(err)
	}
}

// Init DB接続と、マイグレーションを行う。
func Init() {
	initDB()
	// マイグレーション実行
	autoMigration()
}

// Close DBと切断する。
func Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}

// GetDB DB接続情報を返す
func GetDB() *gorm.DB {
	if DB == nil {
		initDB()
	}
	return DB
}

// StartBegin トランザクションを開始する。
func StartBegin() *gorm.DB {
	DB = GetDB()
	tx = DB.Begin()
	return tx
}

// EndRollback トランザクションを終了しロールバックする。
func EndRollback() {
	tx.Rollback()
	tx = nil
}

// EndCommit トランザクションを終了しコミットする。
func EndCommit() {
	tx.Commit()
	tx = nil
}

func autoMigration() {
	fmt.Println("migration")
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Relation{})
}
