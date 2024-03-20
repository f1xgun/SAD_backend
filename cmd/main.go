package main

import (
	"log"
	"sad/internal/app"
)

func main() {
	app, err := app.NewApp()

	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	app.Run()
}
