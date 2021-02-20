#!/usr/bin/env bash
#set -xe

# 脚本路径
sc_dir="$(
  cd "$(dirname "$0")" >/dev/null 2>&1 || exit
  pwd -P
)"

# 去掉路径后缀
rs_path=${sc_dir/grpcs*/grpcs}

# 引入头文件
source $rs_path/bin/libs/headers.sh

# 首先清除已编译文件
ebc_info "首先清除已编译文件..."
#mkdir -p $rs_path/src/rpc/server && rm -rf $rs_path/src/rpc/server/*
rm -rf $rs_path/src/rpc/*/proto/*.pb.go
rm -rf $rs_path/src/rpc/*/proto/*.pb.go

# 编译环境
ENV=${1:-"local"}

# 如果为基础系统
([ "$ENV" != "local" ] && [ "$ENV" != "dev" ] && [ "$ENV" != "test" ] && [ "$ENV" != "beta" ] && [ "$ENV" != "prod" ]) && ebc_error "参数[1: $1]不合法!" && exit

# batch rpc files
rpc_path="${rs_path}/src/rpc"
for rpc_fp in "$rpc_path"/*; do
  rpc_fn=$(basename -- "$rpc_fp")
  echo "[proto->*.pb.go][rpc_fp: $rpc_fp][rpc_fn: $rpc_fn]"
  # compile rpc proto file
  #  protoc -I src/rpc/$rpc_fn src/rpc/$rpc_fn/*.proto --go_out=plugins=grpc:src/rpc/$rpc_fn

  protoc --go_out=src/rpc/$rpc_fn/proto --go-grpc_out=src/rpc/$rpc_fn/proto src/rpc/$rpc_fn/proto/*.proto

  #  protoc -I src/rpc/$rpc_fn --go-grpc_out src/rpc/$rpc_fn/ --go-grpc_opt paths=source_relative src/rpc/$rpc_fn/*.proto
done

#protoc --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true src/rpc/echo/proto/*.proto

#protoc --go_out=src/rpc/rpc_error/proto src/rpc/rpc_error/proto/rpc_error.proto
#protoc --go-grpc_out=src/rpc/rpc_error/proto src/rpc/rpc_error/proto/rpc_error.proto
# 编译文件
#protoc -I src/rpc/proto src/rpc/proto/demo.proto --go-grpc_out=plugins=grpc:src/rpc/server
