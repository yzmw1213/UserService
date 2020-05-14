FROM golang:latest

ENV GO111MODULE=on

WORKDIR /go/src/github.com/yzmw1213/GoMicroApp

COPY . .

EXPOSE 50052