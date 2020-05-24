package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/grpc"
)

func main() {
	start()
}

func start() {
	loadEnv()
	db.Init()
	grpc.NewBlogGrpcServer()
	defer db.Close()
}

func loadEnv() {
	log.Println("Loading env variables")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
