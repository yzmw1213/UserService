#!/bin/sh
DIR=$(cd $(dirname $0); pwd);
SERVER_OUT_DIR="/grpc/blog_grpc";
PROTO_PATH="grpc blog_grpc/blog.proto";

# mode=grpcweb,grpcwebtextで結果が異なる
protoc \
  --proto_path=${PROTO_PATH} \
  --go_out=plugins=grpc:${DIR}${SERVER_OUT_DIR}