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
var blogs []*blog_grpc.Blog

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
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := blog_grpc.NewBlogServiceClient(conn)

	newBlog := &blog_grpc.Blog{
		AuthorId: 12345,
		Title:    "title (edited)",
		Content:  "Content of the first blog (edited)",
	}

	req := &blog_grpc.CreateBlogRequest{
		Blog: newBlog,
	}

	_, err = client.CreateBlog(ctx, req)

	if err != nil {
		t.Fatalf("error occured testing CreateBlog: %v\n", err)
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

	var deleteBlogId int32 = 16

	req := &blog_grpc.DeleteBlogRequest{
		BlogId: deleteBlogId,
	}

	_, err = client.DeleteBlog(ctx, req)

	if err != nil {
		t.Fatalf("error occured testing DeleteBlog: %v\n", err)
	}

	blog := &model.Blog{
		BlogId: deleteBlogId,
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
	db := GetDB()
	if !db.HasTable(model.Blog{}) {
		t.Fatal("db does not have table named 'blogs'")
	}

}

func TestUpdateBlog(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	updateBlog := &model.Blog{
		BlogId:   4,
		AuthorId: 1234567890,
		Title:    "Updated Title",
		Content:  "Updated Content",
	}

	client := blog_grpc.NewBlogServiceClient(conn)
	blogs = append(blogs, &blog_grpc.Blog{
		BlogId:   updateBlog.BlogId,
		AuthorId: updateBlog.AuthorId,
		Title:    updateBlog.Title,
		Content:  updateBlog.Content,
	})
	for _, blog := range blogs {
		req := &blog_grpc.UpdateBlogRequest{
			Blog: blog,
		}
		_, err = client.UpdateBlog(ctx, req)

		if err != nil {
			t.Fatalf("error occured testing UpdateBlog: %v\n", err)
		}
	}

	req := &blog_grpc.ReadBlogRequest{
		BlogId: 4,
	}
	res, err := client.ReadBlog(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.GetBlog().GetAuthorId(), updateBlog.AuthorId, "AuthorId of updated blog is not expectd")
	assert.Equal(t, res.GetBlog().GetTitle(), updateBlog.Title, "Title of updated blog is not expectd")
	assert.Equal(t, res.GetBlog().GetContent(), updateBlog.Content, "Content of updated blog is not expectd")

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
