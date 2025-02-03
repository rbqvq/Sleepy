package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "gitlab.com/CoiaPrant/Sleepy/common"
	"gitlab.com/CoiaPrant/clog"
)

var (
	version = "dev"
	commit  string
)

func main() {
	{
		var logFile string
		var interval int64
		var insecure bool

		flag.StringVar(&logFile, "log", "", "Log file location (default: stdout)")
		flag.StringVar(&rpcAddress, "api", "", "RPC api address")
		flag.StringVar(&rpcSecret, "secret", "", "RPC secret")
		flag.StringVar(&rpcType, "rpc", "google", "RPC type (valid: google)")
		flag.Int64Var(&interval, "interval", 2, "Report state interval")
		flag.StringVar(&reportDeviceName, "device", "", "Report device name (default: follow-system)")
		flag.BoolVar(&insecure, "insecure", false, "Disable SSL/TLS certificate verification")
		flag.BoolVar(clog.DebugFlag(), "debug", false, "Show debug logs")
		help := flag.Bool("h", false, "Show help")
		v := flag.Bool("v", false, "Show version")
		flag.Parse()

		if *help {
			flag.PrintDefaults()
			return
		}

		if *v {
			fmt.Printf("%s-%s", version, commit)
			return
		}

		if logFile != "" {
			if file, err := os.Stat(logFile); err == nil && file.Size() > 1*1024*1024 {
				os.Remove(logFile)
			}

			if logger, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm|os.ModeAppend); err == nil {
				clog.SetOutput(logger)
			}
		}

		if rpcAddress == "" {
			clog.Fatalf("Invalid api argument.")
			return
		}

		if rpcSecret == "" {
			clog.Fatalf("Invalid secret argument.")
			return
		}

		if interval <= 0 {
			clog.Fatalf("Invalid interval argument.")
			return
		}
		reportInterval = time.Duration(interval) * time.Second

		if insecure {
			TLSConfig.InsecureSkipVerify = true
			GoTLSConfig.InsecureSkipVerify = true
		}
	}

	clog.Infof("Sleepy Agent Version: %s-%s", version, commit)

	{
		client, err := NewRPClient()
		if err != nil {
			clog.Fatalf("[gRPC API Client] initial client failed, error: %s", err)
			return
		}
		apiClient = client
		go apiClientDaemon()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP)

	<-sigs
	apiClient.Unregister()
	clog.Print("Exiting")
}
