package module

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var graceful struct {
	sync.Once
	Ctx    context.Context
	Cancel func()
}

func GracefulContext() (context.Context, func()) {
	graceful.Do(func() {
		graceful.Ctx, graceful.Cancel = context.WithCancel(context.Background())

		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Try exit...")
			graceful.Cancel()

			<-c
			log.Println("Force exit")
			os.Exit(0)
		}()
	})
	return graceful.Ctx, graceful.Cancel
}
