package main

import (
	pconfig "github.com/raffaele-pilloni/axxon-test/config"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := pconfig.LoadConfig()
	if err != nil {
		log.Panicf("Error while load config. error: %v", err)
	}
	appHTTPServer, err := http.NewServer(
		config,
	)

	if err != nil {
		log.Panicf("[%s-%s] Http server initialization failed: %v", config.App.AppName, config.App.Env, err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := appHTTPServer.Run(); err != nil {
			log.Panicf("[%s-%s] Http server run failed: %v", config.App.AppName, config.App.Env, err)
		}
	}()

	log.Printf("[%s-%s] Http server started", config.App.AppName, config.App.Env)

	sig := <-sigCh
	log.Printf("[%s-%s] Received signal from os: %s", config.App.AppName, config.App.Env, sig)

	if err := appHTTPServer.Stop(); err != nil {
		log.Panicf("[%s-%s] Http server stop failed: %v", config.App.AppName, config.App.Env, err)
	}

	log.Printf("[%s-%s] Http server stopped", config.App.AppName, config.App.Env)
}
