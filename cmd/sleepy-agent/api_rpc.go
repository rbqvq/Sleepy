package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "gitlab.com/CoiaPrant/Sleepy/common"

	contextHelper "gitlab.com/CoiaPrant/Sleepy/pkg/context"
	"gitlab.com/CoiaPrant/Sleepy/pkg/monitor"
	pb "gitlab.com/CoiaPrant/Sleepy/proto"

	"gitlab.com/CoiaPrant/clog"
	"gitlab.com/go-extension/grpc/credentials"
	"google.golang.org/grpc"
)

type gRPClient struct {
	mutex sync.RWMutex

	client  pb.SleepyClient
	session string

	ctx  context.Context
	done context.CancelCauseFunc
}

func NewGoogleRPClient() (*gRPClient, error) {
	c := new(gRPClient)

	credential := grpc.WithTransportCredentials(credentials.NewTLS(TLSConfig))
	conn, err := grpc.NewClient(rpcAddress,
		credential,
		grpc.WithPerRPCCredentials(&gRPC_Session{
			Session: &c.session,
		}),
	)
	if err != nil {
		clog.Errorf("[gRPC API Client] connect failed, error: %s", err)
		return nil, err
	}

	c.client = pb.NewSleepyClient(conn)
	c.ctx, c.done = contextHelper.WithCancelCause(context.Background())

	err = c.registerDevice()
	if err != nil {
		clog.Errorf("[gRPC API Client] Register failed, error: %s", err)
		return nil, err
	}
	clog.Infof("[gRPC API Client] Register success")

	go c.reportDeviceState()
	return c, nil
}

func (c *gRPClient) registerDevice() error {
	var deviceName string
	if reportDeviceName == "" {
		deviceName = monitor.GetDeviceName()
	}

	response, err := c.client.RegisterDevice(context.Background(), &pb.Device{
		ReportInterval: uint64(reportInterval / time.Second),
		DeviceType:     monitor.GetDeviceType(),
		DevicePlatform: monitor.GetDevicePlatform(),
		DeviceName:     deviceName,
	}, grpc.PerRPCCredentials(&gRPC_Secret{
		Secret: rpcSecret,
	}))
	if err != nil {
		c.done(err)
		return err
	}

	if !response.Ok {
		err := fmt.Errorf("remote message: %s", response.Msg)

		c.done(err)
		return err
	}

	c.session = response.Session
	return nil
}

func (c *gRPClient) reportDeviceState() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			using, appName := monitor.GetDeviceState()
			_, err := c.client.ReportDeviceState(c.ctx, &pb.State{
				Using:   using,
				AppName: appName,
			})
			if err != nil {
				c.done(err)
				clog.Debugf("[gRPC API Client] Report device state failed, error: %s", err)
				return
			}

			clog.Debugf("[gRPC API Client] Reported device state, Using: %t, AppName: %s", using, appName)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *gRPClient) Wait() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	<-c.ctx.Done()
	return c.ctx.Err()
}

func (c *gRPClient) Unregister() (err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.done(nil)

	_, err = c.client.Unregister(context.Background(), &pb.Empty{})
	if err != nil {
		clog.Errorf("[gRPC API Client] Unregister failed, error: %s", err)
		return
	}

	clog.Successf("[gRPC API Client] Unregister successfully")
	return
}
