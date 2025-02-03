package grpc

import (
	"context"

	. "gitlab.com/CoiaPrant/Sleepy/common/server"
	"gitlab.com/CoiaPrant/Sleepy/model"
	"gitlab.com/CoiaPrant/Sleepy/services/device"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CheckSecret(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "failed to get metaData")
	}

	value, ok := md["secret"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metaData not contains secret")
	}
	secret := value[0]

	if secret != Secret {
		return status.Errorf(codes.Unauthenticated, "bad secret")
	}

	return nil
}

func CheckSession(ctx context.Context) (*model.Device, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get metaData")
	}

	value, ok := md["session"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metaData not contains session")
	}
	session_id := value[0]

	device, err := device.GetDevice(session_id)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return device, nil
}

func RemoveDevice(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "failed to get metaData")
	}

	value, ok := md["session"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metaData not contains session")
	}
	session := value[0]

	device.RemoveDevice(session)
	return nil
}
