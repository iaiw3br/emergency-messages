package config

import "github.com/joho/godotenv"

const fileName = ".env"

func Load() error {
	return godotenv.Load(fileName)
}
