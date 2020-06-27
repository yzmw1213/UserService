package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/grpc/blog_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var err error

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	blog_grpc.RegisterBlogServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateBlog(t *testing.T) {
	var createBlogs []*blog_grpc.Blog
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := blog_grpc.NewBlogServiceClient(conn)

	createBlogs = append(createBlogs, &blog_grpc.Blog{
		AuthorId: 1111111111,
		Title:    "Title of the first blog",
		Content:  "Content of the first blog",
	})

	createBlogs = append(createBlogs, &blog_grpc.Blog{
		AuthorId: 222222222,
		Title:    "Title of the secound blog",
		Content:  "Content of the second blog",
	})

	createBlogs = append(createBlogs, &blog_grpc.Blog{
		AuthorId: 333333333,
		Title:    "Title of the third blog",
		Content:  "Content of the third blog",
	})

	createBlogs = append(createBlogs, &blog_grpc.Blog{
		AuthorId: 444444444,
		Title:    "Title of the fourth blog",
		Content:  "Content of the fourth blog",
	})

	for _, blog := range createBlogs {
		req := &blog_grpc.CreateBlogRequest{
			Blog: blog,
		}

		_, err = client.CreateBlog(ctx, req)

		if err != nil {
			t.Fatalf("error occured testing CreateBlog: %v\n", err)
		}
	}
	t.Log("finished TestCreateBlog")
}

func TestDeleteBlog(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := blog_grpc.NewBlogServiceClient(conn)

	var deleteBlogID int32 = 2

	req := &blog_grpc.DeleteBlogRequest{
		BlogId: deleteBlogID,
	}

	_, err = client.DeleteBlog(ctx, req)

	if err != nil {
		t.Fatalf("error occured testing DeleteBlog: %v\n", err)
	}

	blog := &model.Blog{
		BlogID: deleteBlogID,
	}

	if !db.DB.First(&blog).RecordNotFound() {
		t.Fatal("The blog specified was not deleted")
	}
	t.Log("finished TestDeleteBlog")
}

func TestGetDB(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	db := db.GetDB()
	if !db.HasTable(model.Blog{}) {
		t.Fatal("db does not have table named 'blogs'")
	}

}

func TestUpdateBlog(t *testing.T) {
	var updateBlogs []*blog_grpc.Blog

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	updateBlog := &model.Blog{
		BlogID:   3,
		AuthorID: 1234567890,
		Title:    "Updated Title",
		Content:  "Updated Content",
	}

	client := blog_grpc.NewBlogServiceClient(conn)

	updateBlogs = append(updateBlogs, &blog_grpc.Blog{
		BlogId:   updateBlog.BlogID,
		AuthorId: updateBlog.AuthorID,
		Title:    updateBlog.Title,
		Content:  updateBlog.Content,
	})

	for _, blog := range updateBlogs {
		req := &blog_grpc.UpdateBlogRequest{
			Blog: blog,
		}
		_, err = client.UpdateBlog(ctx, req)

		if err != nil {
			t.Fatalf("error occured testing UpdateBlog: %v\n", err)
		}
	}

	req := &blog_grpc.ReadBlogRequest{
		BlogId: 3,
	}

	res, err := client.ReadBlog(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, updateBlog.AuthorID, res.GetBlog().GetAuthorId(), "AuthorId of updated blog is not expectd")
	assert.Equal(t, updateBlog.Title, res.GetBlog().GetTitle(), "Title of updated blog is not expectd")
	assert.Equal(t, updateBlog.Content, res.GetBlog().GetContent(), "Content of updated blog is not expectd")

	t.Log("finished TestUpdateBlog")
}

func TestListBlog(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := blog_grpc.NewBlogServiceClient(conn)

	req := &blog_grpc.ListBlogRequest{}
	stream, err := client.ListBlog(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for {
		_, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happened while listing blog: %v", err)
		}
	}

	t.Log("finished TestListBlog")
}
