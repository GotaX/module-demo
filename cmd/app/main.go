package main

import (
	"log"

	"server/business"
	"server/internal/module"
)

func main() {
	ctx, cancel := module.GracefulContext()
	defer cancel()

	if err := business.Run(ctx); err != nil {
		log.Fatalf("Error exit: %s", err)
	}
}
