
### 开发环境

```bash
brew info protobuf
```

* 安装依赖工具

[tools.go]
```bash
// +build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

* 校验包

```bash
go mod tidy
```

* 执行安装

```bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

* 验证安装版本

```bash
protoc-gen-grpc-gateway -version
protoc-gen-openapiv2 -version
protoc-gen-go-grpc -version
```

### Restful 接口调用 GRpc

grpc-proxy是谷歌的协议缓冲区编译器的一个插件 protoc，它读取protobuf服务定义并生成一个反向代理服务器，该服务器将RESTful HTTP API转换为gRPC。该服务器是根据google.api.http 服务定义中的注释生成的 。

* 安装 `protoc-gen-grpc-gateway` 和 `go` 库

```bash
wget https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.2.0/protoc-gen-grpc-gateway-v2.2.0-darwin-x86_64
```

```bash
# 打开grpc-gateway的github官网 https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.16.0，下载系统对应版本
# 安装grpc-gateway
wget https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v1.16.0/protoc-gen-grpc-gateway-v1.16.0-windows-x86_64.exe -O protoc-gen-grpc-gateway.exe
# 把protoc-gen-grpc-gateway.exe移动到$GOROOT/bin目录下

# 安装go语言使用的grpc-gateway库
wget https://github.com/grpc-ecosystem/grpc-gateway/archive/v1.16.0.zip
# 解压后，把目录改名为grpc-gateway，然后把grpc-gateway移动到$GOPATH/src/github.com/grpc-ecosystem/
# (可选)把grpc-gateway/third_party/googleapis/google/下使用的第三方proto文件复制到$GOROOT/bin/include/google/目录下，主要为了让protoc使用默认路径就可以使用第三方proto文件，如果不复制，需要-I参数来指定目录
```