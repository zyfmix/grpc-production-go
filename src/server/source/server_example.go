package source

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	interceptors2 "grpcs/src/interceptors"
	interceptors "grpcs/src/interceptors/server"
	helloworld "grpcs/src/rpc/helloworld"
	"grpcs/src/tlscert"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type server struct {
	// [mustEmbedUnimplemented*** method appear in grpc-server #3794](https://github.com/grpc/grpc-go/issues/3794)
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("SayHello,Received: %v", in.Name)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Print("SayHello,metadata: ", md)

	authSecurity := ctx.Value("authSecurity")
	if authSecurity, ok := authSecurity.(interceptors.AuthSecurity); ok {
		log.Print("authSecurity.Name:", authSecurity.Name)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname %v", err)
	}
	if hostname != "" {
		grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	}
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func (g *server) SayList(r *helloworld.HelloRequest, stream helloworld.Greeter_SayListServer) error {
	var err error
	fmt.Println("\nSayList receive req: " + r.Name)

	for i := 0; i < 5; i++ {
		err = stream.Send(&helloworld.HelloReply{Message: "hello " + r.Name + fmt.Sprintf(" %d", i)})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}

func (g *server) SayRecord(stream helloworld.Greeter_SayRecordServer) error {
	values := []string{}
	defer func() {
		fmt.Println("\nSayRecord receive req: ", values)
	}()

	for {
		// 阻塞等待接收流数据，当结束时会受到EOF表示结束，当出现错误会返回rpc错误信息
		// 默认的MaxReceiveMessageSize值为1024x1024x4字节，如果有特殊需求可以调整
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				return stream.SendAndClose(&helloworld.HelloReply{
					Message: "hello " + strings.Join(values, ","),
				})
			}
			return err
		}

		values = append(values, resp.Name)
	}

	return nil
}

func (g *server) SayRoute(stream helloworld.Greeter_SayRouteServer) error {
	recValues := []string{}
	sendValues := []string{}

	defer func() {
		fmt.Println("\nSayRoute receive req: ", recValues)
		fmt.Println("SayRoute send req: ", sendValues)
	}()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				return nil
			}
			return err
		}
		recValues = append(recValues, resp.Name)

		err = stream.Send(&helloworld.HelloReply{Message: "hello " + resp.Name})
		if err != nil {
			return err
		}
		sendValues = append(sendValues, resp.Name)
	}

	return nil
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
	s.SetUnaryInterceptors(interceptors2.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(interceptors2.GetDefaultStreamServerInterceptors())
}
