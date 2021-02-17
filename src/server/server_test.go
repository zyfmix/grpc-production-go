package main

import (
	"github.com/stretchr/testify/assert"
	"grpcs/src/interceptors"
	"grpcs/src/server/source"
	"grpcs/src/tlscert"
	"testing"
)

func TestBuildGrpcServer(t *testing.T) {
	builder := &source.GrpcServerBuilder{}
	builder.SetTlsCert(&tlscert.Cert)
	builder.DisableDefaultHealthCheck(true)
	builder.EnableReflection(true)
	builder.SetStreamInterceptors(interceptors.GetDefaultStreamServerInterceptors())
	builder.SetUnaryInterceptors(interceptors.GetDefaultUnaryServerInterceptors())
	server := builder.Build()
	assert.NotNil(t, server)
}
