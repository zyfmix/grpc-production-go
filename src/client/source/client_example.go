package source

import (
	"context"
	"fmt"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"grpcs/src/interceptors"
	helloworld "grpcs/src/rpc/helloworld"
	"grpcs/src/tlscert"
	"io"
	"log"
	"time"
)

func TimeoutLogExample() {
	log.Println("TimeoutLogExample.builder")

	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithStreamInterceptors(interceptors.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(interceptors.GetDefaultUnaryClientInterceptors())

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
	log.Printf("[SayHello:一元调用RPC]###################################################################################################")
	if err = SayHello(client); err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("[SayList:服务端流式RPC]###################################################################################################")
	if err = SayList(client); err != nil {
		log.Printf("err:%v", err)
	}

	log.Printf("[SayRecord:客户端流式RPC]###################################################################################################")
	if err = SayRecord(client); err != nil {
		log.Printf("err:%v", err)
	}

	log.Printf("[SayRoute:双向流式RPC]###################################################################################################")
	if err = SayRoute(client); err != nil {
		log.Printf("err:%v", err)
	}

	log.Println("TimeoutLogExample.ended")

	defer cancel()
}

func SayHello(client helloworld.GreeterClient) error {
	request := &helloworld.HelloRequest{
		Name: "client(mike)",
	}

	helloReply, err := client.SayHello(context.Background(), request)
	if err != nil {
		log.Printf("err:%v", err)
		return err
	}

	fmt.Println("resp:", helloReply.Message)
	return nil
}

func SayList(client helloworld.GreeterClient) error {
	stream, err := client.SayList(context.Background(), &helloworld.HelloRequest{Name: "zhangsan"})
	if err != nil {
		return err
	}

	for {
		// 阻塞等待接收流数据，当结束时会受到EOF表示结束，当出现错误会返回rpc错误信息
		// 默认的MaxReceiveMessageSize值为1024x1024x4字节，如果有特殊需求可以调整
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				break
			}
			return err
		}

		fmt.Println("resp:", resp.Message)
	}

	return nil
}

func SayRecord(client helloworld.GreeterClient) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}

	names := []string{"zhangsan", "lisi", "wangwu"}
	for _, name := range names {
		err := stream.Send(&helloworld.HelloRequest{Name: name})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println("resp:", resp.Message)

	return nil
}

func SayRoute(client helloworld.GreeterClient) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}

	names := []string{"zhangsan", "lisi", "wangwu"}
	for _, name := range names {
		err := stream.Send(&helloworld.HelloRequest{Name: name})
		if err != nil {

			return err
		}

		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Println("resp:", resp.Message)
	}

	time.Sleep(10 * time.Millisecond)
	err = stream.CloseSend()
	if err != nil {
		return err
	}

	return nil
}

func TLSConnExample() {
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithClientTransportCredentials(false, tlscert.CertPool)
	clientBuilder.WithStreamInterceptors(interceptors.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(interceptors.GetDefaultUnaryClientInterceptors())
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
