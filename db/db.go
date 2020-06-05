package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/GoMicroApp/domain/model"

	"github.com/jinzhu/gorm"
)

var (
	DB        *gorm.DB
	blog      model.Blog
	tableName string = "blogs"
)

func initDB() {
	var err error
	DBMS := "mysql"
	DB_ADRESS := os.Getenv("DB_ADRESS")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	PROTOCOL := fmt.Sprintf("tcp(%s)", DB_ADRESS)
	DB_OPTION := "?charset=utf8mb4&parseTime=True&loc=Local"
	CONNECTION := fmt.Sprintf("%s:%s@%s/%s%s", DB_USER, DB_PASSWORD, PROTOCOL, DB_NAME, DB_OPTION)

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

func InsDelUpdOperation(ctx context.Context, op string, postData *model.Blog) error {
	initDB()

	switch op {
	case "insert":
		if err := DB.Create(postData).Error; err != nil {
			return err
		}
	case "update":
		if err := DB.Model(&blog).Updates(postData).Error; err != nil {
			return err
		}
	}
	return nil

}

func Delete(ctx context.Context, postData *model.Blog) error {
	initDB()

	if err := DB.Delete(postData).Error; err != nil {
		return err
	}
	return nil
}

func Read(ctx context.Context, blogId int32) (model.Blog, error) {
	initDB()
	var blog model.Blog

	row := DB.First(&blog, blogId)
	if err := row.Error; err != nil {
		log.Printf("Error happend while Read for blogid: %v\n", blogId)
		return model.Blog{}, err
	}
	DB.Table(tableName).Scan(row)
	return blog, nil
}

func ListAll(ctx context.Context) ([]model.Blog, error) {
	initDB()
	var blog model.Blog
	var blogs []model.Blog
	var rows *sql.Rows

	rows, err := DB.Find(&blogs).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}

	for rows.Next() {
		DB.ScanRows(rows, &blog)
		blogs = append(blogs, blog)
	}
	return blogs, nil
}
