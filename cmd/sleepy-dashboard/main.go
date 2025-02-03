package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "gitlab.com/CoiaPrant/Sleepy/common/server"
	"gitlab.com/CoiaPrant/Sleepy/router"

	"github.com/gin-gonic/gin"
	"gitlab.com/CoiaPrant/clog"
)

var (
	version = "dev"
	commit  string
)

func main() {
	var conf Config
	{
		var config_file, logFile string
		{
			flag.StringVar(&config_file, "config", "config.json", "The config file location")
			flag.StringVar(&logFile, "log", "", "")
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
		}

		file, err := os.ReadFile(config_file)
		if err != nil {
			clog.Fatalf("[Config] Unable to read config file, error: %s", err)
			return
		}

		err = json.Unmarshal(file, &conf)
		if err != nil {
			clog.Fatalf("[Config] Unable to parse config file, error: %s", err)
			return
		}
	}

	{
		// Security
		Secret = conf.Security.Secret
		AllowCORS = conf.Security.AllowCORS

		// System
		if conf.System.DebugMode {
			clog.SetDebug(true)
		}

		// Global
		gin.SetMode(gin.ReleaseMode)
		if clog.IsDebug() {
			gin.SetMode(gin.DebugMode)
		}
	}

	clog.Infof("Sleepy Dashboard Version: %s-%s", version, commit)

	srv := &http.Server{
		Handler:     router.Handler(),
		ErrorLog:    log.New(io.Discard, "", 0),
		IdleTimeout: time.Minute,
		TLSConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(128),
		},
	}
	srv.SetKeepAlivesEnabled(true)

	{
		if conf.Web.Type == "unix" {
			os.Remove(conf.Web.Listen)
		}

		ln, err := (&net.ListenConfig{}).Listen(context.Background(), conf.Web.Type, conf.Web.Listen)
		if err != nil {
			clog.Fatal("[Web] failed to listen, error: ", err)
			return
		}

		if conf.Web.Type == "unix" {
			os.Chmod(conf.Web.Listen, 0777)
		}

		if conf.Web.Cert == "" || conf.Web.Key == "" {
			go srv.Serve(ln)
		} else {
			{
				_, err = tls.LoadX509KeyPair(conf.Web.Cert, conf.Web.Key)
				if err != nil {
					clog.Fatal("[Web] failed to load tls certificate, error: ", err)
					return
				}
			}
			go srv.ServeTLS(ln, conf.Web.Cert, conf.Web.Key)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP)

	<-sigs
	srv.Shutdown(context.Background())
	srv.Close()

	clog.Print("Exiting")
}
