package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yzmw1213/GoMicroApp/grpc/helloworld_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func Run() {
	fmt.Println("Hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	blogServer := &server{}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	helloworld_grpc.RegisterGreeterServer(s, blogServer)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a sgnal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the client")
	lis.Close()
	fmt.Println("End of Program")

}

func (s server) SayHello(ctx context.Context, req *helloworld_grpc.HelloRequest) (*helloworld_grpc.HelloResponse, error) {
	fmt.Printf("request name is : %v\n", req.GetName())
	return &helloworld_grpc.HelloResponse{
		Message: "hello " + req.Name,
	}, nil
}
