package main

import (
	"log"
	"sad/internal/app"
)

func main() {
	application, err := app.NewApp()

	if err != nil {
		log.Fatalf("Failed to init app: %s", err.Error())
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run app: %s", err.Error())
	}

	defer application.CloseDBConnection()
}
