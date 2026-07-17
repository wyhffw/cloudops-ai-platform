package config

import (
	"os"
)

type Config struct {
	Addr    string
	AppName string
	Env     string
}

func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	return Config{
		Addr:    addr,
		AppName: "cloudops-backend",
		Env:     env,
	}
}
