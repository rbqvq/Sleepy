package grpc

import (
	pb "gitlab.com/CoiaPrant/Sleepy/proto"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

var (
	self = grpc.NewServer()
)

func init() {
	pb.RegisterSleepyServer(self, &server{})
}

func Self() *grpc.Server {
	return self
}

func GetServiceInfo() map[string]grpc.ServiceInfo {
	return self.GetServiceInfo()
}

func Serve(lis net.Listener) error {
	return self.Serve(lis)
}

func ServeConn(conn net.Conn) {
	grpc_handleRawConn(self, listenerAddressForServeConn, conn)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.ServeHTTP(w, r)
}

func GracefulStop() {
	self.GracefulStop()
}

func Stop() {
	self.Stop()
}
