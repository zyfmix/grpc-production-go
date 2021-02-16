package testing

import (
	"context"
	"github.com/stretchr/testify/assert"
	helloworld "grpcs/src/rpc/server"
	"grpcs/src/testdata"

	"google.golang.org/grpc"
	"grpcs/src/grpcutils"
	"testing"
)

func startServer() {
	builder := GrpcInProcessingServerBuilder{}
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	server = builder.Build()
	server.RegisterService(func(server *grpc.Server) {
		helloworld.RegisterGreeterServer(server, &testdata.MockedService{})
	})
	server.Start()
}

func TestSayHelloPassingContext(t *testing.T) {
	startServer()
	ctx := context.Background()
	clientBuilder := InProcessingClientBuilder{Server: server}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(ctx)
	clientBuilder.WithOptions(grpc.WithContextDialer(GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost", "8080")

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
	clientBuilder := InProcessingClientBuilder{Server: server}
	clientBuilder.WithInsecure()
	clientBuilder.WithOptions(grpc.WithContextDialer(GetBufDialer(server.GetListener())))
	clientConn, err := clientBuilder.GetConn("localhost", "8080")

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
