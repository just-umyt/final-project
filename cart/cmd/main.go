package main

import (
	"cart/internal/app"
	"log"
)

func main() {
	err := app.RunApp()
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Print("shutdown succes")
	}
}
