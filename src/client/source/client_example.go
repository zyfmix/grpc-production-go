package source

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"grpcs/src/grpcutils"
	helloworld "grpcs/src/rpc/server"
	"grpcs/src/tlscert"
	"log"
	"time"
)

func TimeoutLogExample() {
	log.Println("TimeoutLogExample.builder")

	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithStreamInterceptors(grpcutils.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(grpcutils.GetDefaultUnaryClientInterceptors())

	log.Println("TimeoutLogExample.conn")
	cc, err := clientBuilder.GetConn("localhost:8080")
	defer cc.Close()

	log.Println("TimeoutLogExample.metadata")

	ctx := context.Background()
	md := metadata.Pairs("user", "zhangyafei", "pass", "123456789")
	ctx = metadata.NewOutgoingContext(ctx, md)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	log.Println("TimeoutLogExample.health")

	healthClient := grpc_health_v1.NewHealthClient(cc)
	response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("%v", response)

	log.Println("TimeoutLogExample.sayHello")

	timeout := time.Minute * 1
	ctx, cancel := context.WithTimeout(ctx, timeout)
	client := helloworld.NewGreeterClient(cc)
	request := &helloworld.HelloRequest{
		Name: "client(mike)",
	}
	helloReply, err := client.SayHello(ctx, request)
	if err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("helloReply:%v", helloReply)

	log.Println("TimeoutLogExample.ended")

	defer cancel()
}

func TLSConnExample() {
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithClientTransportCredentials(false, tlscert.CertPool)
	clientBuilder.WithStreamInterceptors(grpcutils.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(grpcutils.GetDefaultUnaryClientInterceptors())
	cc, err := clientBuilder.GetTlsConn("localhost:8080")

	defer cc.Close()
	ctx := context.Background()
	md := metadata.Pairs("user", "zhangyafei", "pass", "123456789")
	ctx = metadata.NewOutgoingContext(ctx, md)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	timeout := time.Minute * 1
	ctx, cancel := context.WithTimeout(ctx, timeout)
	client := helloworld.NewGreeterClient(cc)
	request := &helloworld.HelloRequest{
		Name: "mike",
	}
	healthClient := grpc_health_v1.NewHealthClient(cc)
	response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", response)
	helloReply, err := client.SayHello(ctx, request)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", helloReply)

	defer cancel()
}
