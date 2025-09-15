package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret []byte
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load in")
	}

	jwt_secret := os.Getenv("JWT_SECRET")

	if jwt_secret == "" {
		log.Fatal("missing environment variables")
	}

	return &Config{
		JWTSecret: []byte(jwt_secret),
	}
}
