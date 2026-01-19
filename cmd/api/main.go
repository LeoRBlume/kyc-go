package main

import (
	"fmt"
	"kyc-sim/internal/app"
	"log"
)

func main() {
	app, err := app.BuildApp()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	addr := fmt.Sprintf(":%s", app.Config.HTTPPort)
	log.Printf("listening on %s (env=%s)", addr, app.Config.AppEnv)

	if err := app.Router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
