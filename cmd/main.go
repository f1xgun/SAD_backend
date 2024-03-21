package main

import (
	"log"
	"sad/internal/app"
)

func main() {
	app, err := app.NewApp()

	if err != nil {
		log.Fatalf("Failed to init app: %s", err.Error())
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %s", err.Error())
	}
}
