package grpc

import (
	"net"
	_ "unsafe"

	"google.golang.org/grpc"
)

const (
	listenerAddressForServeConn = "listenerAddressForServeConn"
)

//go:linkname grpc_handleRawConn google.golang.org/grpc.(*Server).handleRawConn
func grpc_handleRawConn(server *grpc.Server, lisAddr string, rawConn net.Conn)