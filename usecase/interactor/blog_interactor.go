package interactor

import (
	"context"
	"fmt"

	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/usecase/repository"
)

var err error

type BlogInteractor struct {
	repository.BlogRepository
}

func NewBlogInteractor() *BlogInteractor {
	return &BlogInteractor{}
}

func (b *BlogInteractor) Create(inputBlog *model.Blog) error {
	if err := db.InsDelUpdOperation(context.Background(), "insert", inputBlog); err != nil {
		return err
	}
	return nil
}

func (b *BlogInteractor) CreateBlog(postData *model.Blog) error {

	if err = b.Create(postData); err != nil {
		return err
	}

	return nil
}

func (b *BlogInteractor) Delete(inputBlog *model.Blog) error {
	if err := db.Delete(context.Background(), inputBlog); err != nil {
		return err
	}
	return nil
}

func (b *BlogInteractor) DeleteBlog(postData *model.Blog) error {
	var err error
	if err = b.Delete(postData); err != nil {
		return err
	}
	return nil
}
func (b *BlogInteractor) List() ([]model.Blog, error) {
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

func (b *BlogInteractor) ListBlog() ([]model.Blog, error) {
	return b.List()
}

func (b *BlogInteractor) Update(inputBlog *model.Blog) error {
	if err := db.InsDelUpdOperation(context.Background(), "update", inputBlog); err != nil {
		return err
	}
	return nil
}

func (b *BlogInteractor) UpdateBlog(postData *model.Blog) error {
	var err error

	if err = b.Update(postData); err != nil {
		return err
	}

	return nil
}

func (b *BlogInteractor) Read(blogId int32) (model.Blog, error) {
	return db.Read(context.Background(), blogId)
}

func (b *BlogInteractor) ReadBlog(blogId int32) (model.Blog, error) {
	return b.Read(blogId)
}
