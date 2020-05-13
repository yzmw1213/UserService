package main

import (
	"github.com/yzmw1213/GoMicroApp/server/grpc"
)

func main() {
	start()
}

func start() {
	s := grpc.Greeter{}
	s.Run()
}
