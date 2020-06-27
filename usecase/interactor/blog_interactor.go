package interactor

import (
	"context"
	"fmt"
	"log"

	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/usecase/repository"
)

var (
	err  error
	blog model.Blog
)

// BlogInteractor 投稿サービスを提供するメソッド群
type BlogInteractor struct {
	repository.BlogRepository
}

// CreateBlog 投稿1件を作成
func (b *BlogInteractor) CreateBlog(postData *model.Blog) error {
	DB := db.GetDB()
	if err := DB.Create(postData).Error; err != nil {
		return err
	}

	return nil
}

// DeleteBlog 投稿1件を削除
func (b *BlogInteractor) DeleteBlog(postData *model.Blog) error {
	DB := db.GetDB()
	if err := DB.Delete(postData).Error; err != nil {
		return err
	}
	return nil
}

// ListBlog 投稿を全件取得
func (b *BlogInteractor) ListBlog() ([]model.Blog, error) {
	var blogList []model.Blog
	rows, err := db.ListAll(context.Background())
	if err != nil {
		fmt.Println("Error happened")
		return []model.Blog{}, err
	}
	for _, row := range rows {
		blogList = append(blogList, row)
	}

	return blogList, nil
}

// UpdateBlog 投稿を更新する
func (b *BlogInteractor) UpdateBlog(postData *model.Blog) error {
	DB := db.GetDB()
	if err := DB.Model(&blog).Updates(postData).Error; err != nil {
		return err
	}

	return nil
}

// ReadBlog IDを元に投稿を1件取得する
func (b *BlogInteractor) ReadBlog(blogID int32) (model.Blog, error) {
	DB := db.GetDB()
	row := DB.First(&blog, blogID)
	if err := row.Error; err != nil {
		log.Printf("Error happend while Read for blogiD: %v\n", blogID)
		return model.Blog{}, err
	}
	DB.Table(db.TableName).Scan(row)
	return blog, nil
}
