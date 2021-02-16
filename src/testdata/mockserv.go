package testdata

import (
	"context"
	helloworld "grpcs/src/rpc/server"
)

type MockedService struct {
	helloworld.UnimplementedGreeterServer
}

func (s *MockedService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "This is a mocked service " + in.Name}, nil
}
