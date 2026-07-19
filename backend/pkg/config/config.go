package config

import (
	"os"
	"time"
)

type Config struct {
	Addr     string
	AppName  string
	Env      string
	Version  string
	JWTSecret string
	JWTExpire time.Duration
	AdminUser string
	AdminPass string
	KubeConfig string
}

func Load() Config {
	addr := envOr("ADDR", ":8080")
	env := envOr("APP_ENV", "dev")
	expire := 24 * time.Hour

	return Config{
		Addr:       addr,
		AppName:    "cloudops-backend",
		Env:        env,
		Version:    "0.3.0",
		JWTSecret:  envOr("JWT_SECRET", "cloudops-dev-secret-change-me"),
		JWTExpire:  expire,
		AdminUser:  envOr("ADMIN_USER", "admin"),
		AdminPass:  envOr("ADMIN_PASSWORD", "admin123"),
		KubeConfig: os.Getenv("KUBECONFIG"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
