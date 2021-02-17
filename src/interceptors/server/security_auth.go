package interceptors

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthentication() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(securityContextHandle)
}

func StreamAuthentication() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(securityContextHandle)
}

func securityContextHandle(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	user, ok := md["user"]
	pass, ok := md["pass"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied.")
	}

	if user[0] != "zhangyafeing" || pass[0] != "123456789" {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied...")
	}

	newCtx := context.WithValue(ctx, "authSecurity", AuthSecurity{"foo-s"})
	return newCtx, nil
}

type AuthSecurity struct {
	Name string
}
