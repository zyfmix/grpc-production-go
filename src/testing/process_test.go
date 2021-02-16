package testing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"grpcs/src/grpcutils"
	helloworld "grpcs/src/rpc/server"
	"grpcs/src/testdata"
	"testing"
)

var server GrpcInProcessingServer

func serverStart() {
	builder := GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}

//TestSayHello will test the HelloWorld service using A in memory data transfer instead of the normal networking
func TestSayHello(t *testing.T) {
	serverStart()
	ctx := context.Background()
	clientConn, err := GetInProcessingClientConn(ctx, server.GetListener(), []grpc.DialOption{})
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
