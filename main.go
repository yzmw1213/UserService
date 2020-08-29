package main

import (
	"github.com/yzmw1213/UserService/db"
	"github.com/yzmw1213/UserService/grpc"
)

func main() {
	start()
}

func start() {
	db.Init()
	grpc.NewUserGrpcServer()
	defer db.Close()
}
