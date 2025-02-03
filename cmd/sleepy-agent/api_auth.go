package main

import (
	"context"
)

type gRPC_Secret struct {
	Secret string
}

func (c *gRPC_Secret) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"secret": c.Secret,
	}, nil
}

func (*gRPC_Secret) RequireTransportSecurity() bool {
	return false
}

type gRPC_Session struct {
	Session *string
}

func (c *gRPC_Session) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"session": *c.Session,
	}, nil
}

func (*gRPC_Session) RequireTransportSecurity() bool {
	return false
}
