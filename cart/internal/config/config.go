package config

import (
	"github.com/joho/godotenv"
)

func LoadConfig(path string) error {
	return godotenv.Load(path)
}
