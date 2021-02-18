package main

import (
	echo "grpcs/src/rpc/echo/proto"
	"grpcs/src/rpc/echo/server/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[echo] ")
}

func main() {
	port := ":38080"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	srv := grpc.NewServer()
	echo.RegisterEchoServiceServer(srv, &service.EchoService{})
	log.Printf("start server on port%s\n", port)
	if err := srv.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
