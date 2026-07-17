package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/api"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api.RegisterRoutes(r, cfg)

	log.Printf("starting %s on %s (env=%s)", cfg.AppName, cfg.Addr, cfg.Env)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
