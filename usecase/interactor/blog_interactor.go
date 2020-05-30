package interactor

import (
	"context"

	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/usecase/repository"
)

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
	var err error

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
