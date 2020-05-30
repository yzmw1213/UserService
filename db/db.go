package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/GoMicroApp/domain/model"

	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

func initDB() {
	var err error
	DBMS := "mysql"
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

func InsDelUpdOperation(ctx context.Context, op string, postData *model.Blog) error {
	initDB()

	switch op {
	case "insert":
		if err := DB.Create(postData).Error; err != nil {
			return err
		}
	}
	return nil

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
