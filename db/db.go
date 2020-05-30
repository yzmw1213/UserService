package db

import (
	"context"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/GoMicroApp/domain/model"

	"github.com/jinzhu/gorm"
)

var (
	DB   *gorm.DB
	blog model.Blog
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
	case "update":
		if err := DB.Model(&blog).Updates(postData).Error; err != nil {
			return err
		}
	}
	return nil

}
