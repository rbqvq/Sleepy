package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/CoiaPrant/Sleepy/router/grpc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Handler() http.Handler {
	return h2c.NewHandler(&handler{gin: NewServer()}, &http2.Server{})
}

type handler struct {
	gin *gin.Engine
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.RemoteAddr, ":") {
		r.RemoteAddr = "127.0.0.1:0"
	}

	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {
		grpc.ServeHTTP(w, r)
	} else {
		h.gin.ServeHTTP(w, r)
	}
}
