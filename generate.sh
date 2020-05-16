#!/bin/sh
DIR=$(cd $(dirname $0); pwd);
SERVER_OUT_DIR="/grpc/helloworld_grpc";
PROTO_PATH="grpc helloworld_grpc/helloworld.proto";

# mode=grpcweb,grpcwebtextで結果が異なる
protoc \
  --proto_path=${PROTO_PATH} \
  --go_out=plugins=grpc:${DIR}${SERVER_OUT_DIR}