#!/usr/bin/env bash
set -xe

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
mkdir -p $rs_path/src/rpc/server && rm -rf $rs_path/src/rpc/server/*

# 编译环境
ENV=${1:-"local"}

# 如果为基础系统
([ "$ENV" != "local" ] && [ "$ENV" != "dev" ] && [ "$ENV" != "test" ] && [ "$ENV" != "beta" ] && [ "$ENV" != "prod" ]) && ebc_error "参数[1: $1]不合法!" && exit

# batch rpc files
rpc_path="${rs_path}/src/rpc"
for rpc_fp in "$rpc_path"/proto/*; do
  rpc_fn=$(basename -- "$rpc_fp")
  echo "[proto->server][rpc_fp: $rpc_fp][rpc_fn: $rpc_fn]"
  # compile rpc proto file
  protoc -I src/rpc/proto src/rpc/proto/$rpc_fn --go_out=plugins=grpc:src/rpc/server
done

# 编译文件
#protoc -I src/rpc/proto src/rpc/proto/demo.proto --go_out=plugins=grpc:src/rpc/server
