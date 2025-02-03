//go:build dev

package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"gitlab.com/CoiaPrant/clog"
)

func init() {
	go func() {
		ln, err := net.ListenTCP("tcp", nil)
		if err != nil {
			clog.Errorf("[Debug] Failed to listen pprof, error: %s", err)
			return
		}
		defer ln.Close()

		clog.Infof("[Debug] Listening pprof on %s", ln.Addr())

		err = http.Serve(ln, nil)
		if err != nil {
			clog.Errorf("[Debug] Failed to serve pprof, error: %s", err)
			return
		}
	}()
}
