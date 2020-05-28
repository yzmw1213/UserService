package grpc

import (
	"context"
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
