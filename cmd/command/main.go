package main

import (
	"context"
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/command"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	configs, err := pconfigs.LoadConfigs()
	if err != nil {
		log.Panicf("Error while load configs. error: %v", err)
	}
	commandDispatcher, err := command.NewDispatcher(
		configs,
	)

	if err != nil {
		log.Panicf("Command dispatcher initialization failed: %s", err)
	}

	if len(os.Args) < 2 {
		log.Panicf("Command name must be defined: %s", err)
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
			log.Panicf("Command %s run failed: %s", commandName, err)
		}
	}()

	log.Printf("Command %s started", commandName)

	sig := <-sigCh
	log.Printf("Received signal from os: %s", sig)

	cancelCtx()
	wg.Wait()

	log.Printf("Command %s stopped", commandName)
}
