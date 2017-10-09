package main

import (
	"log"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		select {
		case <-sigs:
			cancel()
			log.Printf("caught signal, stopping...")
		}
	}()
	if runner, err := setupRunnerWithConfiguration(ctx); nil != err {
		log.Printf("failed on starting up runner: %v", err)
	} else {
		runner.Run()
	}
	log.Print("stopped.")
}