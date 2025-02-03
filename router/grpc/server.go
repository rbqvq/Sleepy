package grpc

import (
	"context"
	"time"

	pb "gitlab.com/CoiaPrant/Sleepy/proto"
	"gitlab.com/CoiaPrant/Sleepy/services/device"
)

type server struct {
	pb.UnimplementedSleepyServer
}

func (s *server) RegisterDevice(ctx context.Context, dev *pb.Device) (*pb.RegisterResponse, error) {
	err := CheckSecret(ctx)
	if err != nil {
		return nil, err
	}

	session_id := device.AddDevice(dev)
	return &pb.RegisterResponse{
		Ok:      true,
		Session: session_id,
		Msg:     "Register successfully",
	}, nil
}

func (s *server) Unregister(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, RemoveDevice(ctx)
}

func (s *server) ReportDeviceState(ctx context.Context, state *pb.State) (*pb.Empty, error) {
	device, err := CheckSession(ctx)
	if err != nil {
		return nil, err
	}

	device.AppName = state.AppName
	device.Using = state.Using
	device.Timestamp = time.Now().Unix()
	return &pb.Empty{}, nil
}
