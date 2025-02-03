package main

import (
	//. "PortForwardGo/common"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	apiClient *gRPClient

	rpcType    string = "google"
	rpcAddress string
	rpcSecret  string

	reportInterval   time.Duration = 2 * time.Second
	reportDeviceName string
)

func apiClientDaemon() {
	for {
		err := apiClient.Wait()
		if errors.Is(err, context.Canceled) {
			return
		}

		time.Sleep(5 * time.Second)

		client, err := NewRPClient()
		if err != nil {
			continue
		}

		apiClient = client
	}
}

func NewRPClient() (*gRPClient, error) {
	switch rpcType {
	case "", "google":
		return NewGoogleRPClient()
	default:
		return nil, fmt.Errorf("unknown rpc mode: %s", rpcType)
	}
}
