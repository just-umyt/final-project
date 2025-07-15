package main

import (
	"flag"
	"log"

	"metrics-consumer/internal/app"
)

var (
	ErrLoadEnv = "error loading .env file: %v"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "prod", `There are 2 env: 1 - "prod", 2 - "local"`)
	flag.Parse()

	env = ".env." + env

	err := app.RunApp(env)
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Print("shutdown succes")
	}
}
