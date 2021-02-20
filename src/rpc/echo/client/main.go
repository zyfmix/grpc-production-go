package main

import (
	"context"
	echo "grpcs/src/rpc/echo/proto"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[echo] ")
}

func main() {
	target := "localhost:38080"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	client := echo.NewEchoServiceClient(conn)

	// echo msg
	msg := "says"
	if len(os.Args) > 1 {
		msg = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Echo(ctx, &echo.EchoRequest{Message: msg})
	if err != nil {
		log.Println(err)
	}
	log.Println(r.GetMessage())
}
