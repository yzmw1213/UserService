package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

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

	client := blog_grpc.NewBlogServiceClient(conn)
	blogs = append(blogs, &blog_grpc.Blog{
		BlogId:   4,
		AuthorId: 4444444,
		Title:    "Title (Reading)",
		Content:  "Content (Reading)",
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
	if res.GetBlog().GetAuthorId() != 4444444 || res.GetBlog().GetTitle() != "Title (Reading)" || res.GetBlog().GetContent() != "Content (Reading)" {
		t.Fatal("Result of TestReadBlog was unexpected!")
	}

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
