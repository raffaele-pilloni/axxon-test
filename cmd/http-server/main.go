package main

import (
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configs, err := pconfigs.LoadConfigs()
	if err != nil {
		log.Panicf("Error while load configs. error: %v", err)
	}
	appHTTPServer, err := http.NewServer(
		configs,
	)

	if err != nil {
		log.Panicf("[%s-%s] Http server initialization failed: %v", configs.App.AppName, configs.App.Env, err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := appHTTPServer.Run(); err != nil {
			log.Panicf("[%s-%s] Http server run failed: %v", configs.App.AppName, configs.App.Env, err)
		}
	}()

	log.Printf("[%s-%s] Http server started", configs.App.AppName, configs.App.Env)

	sig := <-sigCh
	log.Printf("[%s-%s] Received signal from os: %s", configs.App.AppName, configs.App.Env, sig)

	if err := appHTTPServer.Stop(); err != nil {
		log.Panicf("[%s-%s] Http server stop failed: %v", configs.App.AppName, configs.App.Env, err)
	}

	log.Printf("[%s-%s] Http server stopped", configs.App.AppName, configs.App.Env)
}
