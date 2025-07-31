package main

import (
	"flag"
	"stocks/internal/app"
	myZap "stocks/pkg/zap"
)

const (
	logFilePath = "././app.log"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "prod", `There are 2 env: 1 - "prod", 2 - "local"`)
	flag.Parse()

	env = ".env." + env

	logger, cleanup, err := myZap.NewLogger(logFilePath)
	if err != nil {
		logger.Errorf("error NewLogger: %v", err)
	}
	defer cleanup()

	err = app.RunApp(env, logger)
	if err != nil {
		logger.Fatalf("error: %v", err)
	} else {
		logger.Info("shutdown succes")
	}
}
