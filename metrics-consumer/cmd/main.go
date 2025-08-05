package main

import (
	"flag"
	"fmt"
	"log"

	"metrics-consumer/internal/app"
	"metrics-consumer/internal/config"
)

const (
	ErrLoadEnv = "error loading .env file: %v"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "prod", `There are 2 env: 1 - "prod", 2 - "local"`)
	flag.Parse()

	env = ".env." + env
	if err := config.LoadConfig(env); err != nil {
		err = fmt.Errorf(ErrLoadEnv, err)
		log.Fatalf(ErrLoadEnv, err)
	}

	err := app.RunApp()
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Print("shutdown succes")
	}
}
