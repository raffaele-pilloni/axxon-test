package main

import (
	pconfig "github.com/raffaele-pilloni/axxon-test/config"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http"
	clog "github.com/raffaele-pilloni/axxon-test/internal/log"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := pconfig.LoadConfig(false)
	if err != nil {
		log.Panicf("Load configuration failed. error: %v", err)
	}

	if err := clog.InitLogConfiguration(
		config.App.ProjectDir,
		config.App.Env,
		config.App.AppName,
		config.App.ServiceName,
		os.Stdout,
	); err != nil {
		log.Panicf("Init log configuration failed. error: %v", err)
	}

	appHTTPServer, err := http.NewServer(
		config,
	)

	if err != nil {
		log.Panicf("Http server initialization failed: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := appHTTPServer.Run(); err != nil {
			log.Panicf("Http server run failed: %v", err)
		}
	}()

	log.Printf("Http server started")

	sig := <-sigCh
	log.Printf("Received signal from os: %s", sig)

	if err := appHTTPServer.Stop(); err != nil {
		log.Panicf("Http server stop failed: %v", err)
	}

	log.Printf("Http server stopped")
}
