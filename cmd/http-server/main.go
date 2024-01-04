package main

import (
	"github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/**
 * Main.
 */
func main() {
	configs, err := configs.LoadConfigs()
	if err != nil {
		log.Panicf("Error while load configs. error: %v", err)
	}
	appHttpServer, err := http.NewServer(
		configs,
	)

	if err != nil {
		log.Panicf("app http server initialization failed: %s", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := appHttpServer.Run(); err != nil {
			log.Panicf("app http server run failed: %s", err)
		}
	}()

	log.Print("app http server started")

	sig := <-sigCh
	log.Printf("received signal from os: %s", sig)

	if err := appHttpServer.Stop(); err != nil {
		log.Panicf("app http server stop failed: %s", err)
	}

	log.Print("app http server stopped")
}