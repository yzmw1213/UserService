#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
OUT_DIR="/grpc/user_grpc";
PROTO_FILE="user.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${OUT_DIR}" \
  -I".${OUT_DIR}" "${PROTO_FILE}"
