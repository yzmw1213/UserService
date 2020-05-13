package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yzmw1213/GoMicroApp/server/helloworld"
	pb "github.com/yzmw1213/GoMicroApp/server/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Greeter struct{}

//main make server implement interface
func (g Greeter) Run() {
	fmt.Println("Hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	// echopb.RegisterEchoServiceServer(s, &server{})
	pb.RegisterGreeterServer(s, g)

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

func (g Greeter) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	fmt.Printf("request name is : %v\n", req.GetName())
	return &helloworld.HelloResponse{
		Message: "hello " + req.Name,
	}, nil
}
