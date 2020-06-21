#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
OUT_DIR="/grpc/blog_grpc";
PROTO_FILE="blog.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${OUT_DIR}" \
  -I".${OUT_DIR}" "${PROTO_FILE}"
