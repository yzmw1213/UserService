#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
USER_OUT_DIR="/grpc/userservice";
POST_OUT_DIR="/grpc/postservice";
USER_PROTO_FILE="user.proto";
POST_PROTO_FILE="post.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${USER_OUT_DIR}" \
  -I".${USER_OUT_DIR}" "${USER_PROTO_FILE}"

protoc \
  --go_out=plugins=grpc:"${DIR}${POST_OUT_DIR}" \
  -I".${POST_OUT_DIR}" "${POST_PROTO_FILE}"
