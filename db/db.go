package db

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/grpc/blog_grpc"

	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

func initDB() {
	var err error
	// env, err := godotenv.Read()
	DBMS := "mysql"
	// DB_USER := env["DB_USER"]
	// DB_PASS := env["DB_PASS"]
	// DB_ADRESS := env["DB_ADRESS"]
	// DB_NAME := env["DB_NAME"]
	// PROTOCOL := fmt.Sprintf("tcp(%s)", DB_ADRESS)
	// DB_NAME := os.Getenv("DB_NAME")
	// DB_OPTION := "?charset=utf8mb4&parseTime=True&loc=Local"
	// CONNECTION := fmt.Sprintf("%s:%s@%s/%s%s", DB_USER, DB_PASS, PROTOCOL, DB_NAME, DB_OPTION)
	CONNECTION := "yzmw1213:root@tcp(localhost)/db?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(DBMS, CONNECTION)
	if err != nil {
		panic(err)
	}
}

func Init() {
	initDB()
	// マイグレーション実行
	autoMigration()
}

func Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	initDB()
	return DB
}

func autoMigration() {
	fmt.Println("migration")
	err := DB.AutoMigrate(&model.Blog{}).Error
	if err != nil {
		panic(err)
	}
}

func InsDelUpdOperation(ctx context.Context, op string, postData *blog_grpc.Blog) error {
	blog := &model.Blog{
		AuthorId: postData.GetAuthorId(),
		Title:    postData.GetTitle(),
		Content:  postData.GetContent(),
	}
	log.Printf("inserting blog AuthorId: %v\n", blog.AuthorId)
	log.Printf("inserting blog Title: %v\n", blog.Title)
	log.Printf("inserting blog Content: %v\n", blog.Content)
	log.Printf("op :%v", op)
	initDB()

	switch op {
	case "insert":
		if err := DB.Create(blog).Error; err != nil {
			return err
		}
	}
	return nil

}

func SelectFirst() model.Blog {
	blog := model.Blog{}
	GetDB().First(blog)
	return blog
}
