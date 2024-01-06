package main

import (
	"context"
	pconfig "github.com/raffaele-pilloni/axxon-test/config"
	"github.com/raffaele-pilloni/axxon-test/internal/app/command"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config, err := pconfig.LoadConfig()
	if err != nil {
		log.Panicf("Error while load config. error: %v", err)
	}
	commandDispatcher, err := command.NewDispatcher(
		config,
	)

	if err != nil {
		log.Panicf("[%s-%s] Command dispatcher initialization failed: %s", config.App.AppName, config.App.Env, err)
	}

	if len(os.Args) < 2 {
		log.Panicf("[%s-%s] Command name must be defined: %s", config.App.AppName, config.App.Env, err)
	}

	commandName := os.Args[1]
	args := os.Args[:2]

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := commandDispatcher.Run(ctx, commandName, args); err != nil {
			log.Panicf("[%s-%s] Command %s run failed: %v", config.App.AppName, config.App.Env, commandName, err)
		}
	}()

	log.Printf("[%s-%s] Command %s started", config.App.AppName, config.App.Env, commandName)

	sig := <-sigCh
	log.Printf("[%s-%s] Received signal from os: %s", config.App.AppName, config.App.Env, sig)

	cancelCtx()
	wg.Wait()

	log.Printf("[%s-%s] Command %s stopped", config.App.AppName, config.App.Env, commandName)
}
