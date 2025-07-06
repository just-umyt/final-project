package main

import (
	"cart/internal/app"
	"flag"
	"log"
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
