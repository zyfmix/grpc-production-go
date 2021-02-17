package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	source2 "grpcs/src/client/source"
	"grpcs/src/grpcutils"
	helloworld "grpcs/src/rpc/server"
	"grpcs/src/server/source"
	"grpcs/src/testdata"
	gtest "grpcs/src/testing"
	"grpcs/src/tlscert"
	"testing"
)

var server gtest.GrpcInProcessingServer

func startServer() {
	builder := gtest.GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}
func startServerWithTLS() source.GrpcServer {
	builder := source.GrpcServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	builder.SetTlsCert(&tlscert.Cert)
	svr := builder.Build()
	svr.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	svr.Start("localhost:8989")
	return svr
}

func TestSayHelloPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := source2.GrpcConnBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(ctx)
	clientBuilder.WithOptions(grpc.WithContextDialer(gtest.GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost:8080")

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}

func TestSayHelloNotPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := source2.GrpcConnBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithOptions(grpc.WithContextDialer(gtest.GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost:8080")

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}

func TestTLSConnWithCert(t *testing.T) {
	serverWithTLS := startServerWithTLS()
	defer serverWithTLS.GetListener().Close()

	ctx := context.Background()
	clientBuilder := source2.GrpcConnBuilder{}
	clientBuilder.WithContext(ctx)
	clientBuilder.WithBlock()
	clientBuilder.WithClientTransportCredentials(false, tlscert.CertPool)
	clientConn, _ := clientBuilder.GetTlsConn("localhost:8989")
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Message, "This is a mocked service test")
}

func TestTLSConnWithInsecure(t *testing.T) {
	serverWithTLS := startServerWithTLS()
	defer serverWithTLS.GetListener().Close()

	ctx := context.Background()
	clientBuilder := source2.GrpcConnBuilder{}
	clientBuilder.WithContext(ctx)
	clientBuilder.WithBlock()
	clientBuilder.WithClientTransportCredentials(true, nil)
	clientConn, _ := clientBuilder.GetTlsConn("localhost:8989")
	defer clientConn.Close()
	client := helloworld.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
