package source

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpcs/src/grpcutils"
	helloworld "grpcs/src/rpc/server"
	"grpcs/src/tlscert"
	"log"
	"os"
)

type server struct {
	// [mustEmbedUnimplemented*** method appear in grpc-server #3794](https://github.com/grpc/grpc-go/issues/3794)
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("SayHello,Received: %v", in.Name)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Print("SayHello,", md)
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname %v", err)
	}
	if hostname != "" {
		grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	}
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func ServerInitialization() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	builder := GrpcServerBuilder{}
	addInterceptors(&builder)
	builder.EnableReflection(true)
	s := builder.Build()
	s.RegisterService(serviceRegister)
	err := s.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func ServerInitializationWithTLS() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	builder := GrpcServerBuilder{}
	addInterceptors(&builder)
	builder.EnableReflection(true)

	// setter tls cert
	builder.SetTlsCert(&tlscert.Cert)

	s := builder.Build()
	s.RegisterService(serviceRegister)
	err := s.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func serviceRegister(sv *grpc.Server) {
	helloworld.RegisterGreeterServer(sv, &server{})
}

func addInterceptors(s *GrpcServerBuilder) {
	s.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(grpcutils.GetDefaultStreamServerInterceptors())
}
