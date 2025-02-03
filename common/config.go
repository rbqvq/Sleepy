package common

import (
	"crypto/tls"

	"gitlab.com/go-extension/tlsfinger"
)

var (
	TLSConfig   = tlsfinger.TLSConfig(tlsfinger.Finger_Edge_106, &tlsfinger.Config{})
	GoTLSConfig = &tls.Config{}
)
