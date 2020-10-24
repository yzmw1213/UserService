package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yzmw1213/UserService/grpc/user_grpc"
	"github.com/yzmw1213/UserService/usecase/interactor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	Usecase interactor.UserInteractor
}

// NewUserGrpcServer gRPCサーバー起動
func NewUserGrpcServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := &server{}

	s := makeServer()

	user_grpc.RegisterUserServiceServer(s, server)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("main grpc server has started")

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

func makeServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(transmitStatusInterceptor)),
		// grpc_auth.UnaryServerInterceptor(authorization.AuthFunc),
		// ),
	)
	return s
}

// func (s *server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
// 	// HealthCheck がコールされた場合は認証のチェックを skip する
// 	if fullMethodName == "/hoge_proto.HogeService/HealthCheck" {
// 		return ctx, nil
// 	}
// }
