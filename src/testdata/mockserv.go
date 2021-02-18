package testdata

import (
	"context"
	"grpcs/src/rpc/helloworld"
	"grpcs/src/rpc/helloworld/proto"
)

type MockedService struct {
	proto.UnimplementedGreeterServer
}

func (s *MockedService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "This is a mocked service " + in.Name}, nil
}
