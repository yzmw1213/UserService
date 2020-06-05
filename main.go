package main

import (
	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/grpc"
)

func main() {
	start()
}

func start() {
	db.Init()
	grpc.NewBlogGrpcServer()
	defer db.Close()
}
