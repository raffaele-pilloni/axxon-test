package main

import (
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
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
	configs, err := pconfigs.LoadConfigs()
	if err != nil {
		log.Panicf("Error while load configs. error: %v", err)
	}
	appHTTPServer, err := http.NewServer(
		configs,
	)

	if err != nil {
		log.Panicf("App http server initialization failed: %s", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := appHTTPServer.Run(); err != nil {
			log.Panicf("App http server run failed: %s", err)
		}
	}()

	log.Print("App http server started")

	sig := <-sigCh
	log.Printf("Received signal from os: %s", sig)

	if err := appHTTPServer.Stop(); err != nil {
		log.Panicf("App http server stop failed: %s", err)
	}

	log.Print("App http server stopped")
}
