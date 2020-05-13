#!/bin/sh
DIR=$(cd $(dirname $0); pwd);
SERVER_OUT_DIR="/server/helloworld";
PROTO_PATH="protoc helloworld.proto";

if [ -e $DIR$SERVER_OUT_DIR ]; then
  break
else
  mkdir $DIR$SERVER_OUT_DIR
fi

# mode=grpcweb,grpcwebtextで結果が異なる
protoc \
  --proto_path=${PROTO_PATH} \
  --go_out=plugins=grpc:${DIR}${SERVER_OUT_DIR}