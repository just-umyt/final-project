package main

import (
	"log"
	"stocks/internal/app"
)

func main() {
	err := app.RunApp()
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Print("shutdown succes")
	}
}
