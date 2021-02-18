package service

import (
	"context"
	echo "grpcs/src/rpc/echo/proto"
	"log"
)

type EchoService struct {
	echo.UnimplementedEchoServiceServer
}

func (s *EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	log.Println("Echo,Msg: ", req.Message)
	return &echo.EchoResponse{Message: req.GetMessage()}, nil
}
