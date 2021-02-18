package testing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"grpcs/src/interceptors"
	"grpcs/src/rpc/helloworld"
	"grpcs/src/rpc/helloworld/proto"
	"grpcs/src/testdata"
	"testing"
)

var server GrpcInProcessingServer

func serverStart() {
	builder := GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(interceptors.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		proto.RegisterGreeterServer(server, &testdata.MockedService{})
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
	client := proto.NewGreeterClient(clientConn)
	request := &helloworld.HelloRequest{Name: "test"}
	resp, err := client.SayHello(ctx, request)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	server.Cleanup()
	clientConn.Close()
	assert.Equal(t, resp.Message, "This is a mocked service test")
}
