package source

import (
	"github.com/stretchr/testify/assert"
	"grpcs/src/grpcutils"
	"grpcs/src/tlscert"
	"testing"
)

func TestBuildGrpcServer(t *testing.T) {
	builder := &GrpcServerBuilder{}
	builder.SetTlsCert(&tlscert.Cert)
	builder.DisableDefaultHealthCheck(true)
	builder.EnableReflection(true)
	builder.SetStreamInterceptors(grpcutils.GetDefaultStreamServerInterceptors())
	builder.SetUnaryInterceptors(grpcutils.GetDefaultUnaryServerInterceptors())
	server := builder.Build()
	assert.NotNil(t, server)
}
