package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/api"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/k8s"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	client, err := k8s.NewClient(cfg.KubeConfig)
	if err != nil {
		log.Printf("warning: kubernetes client init failed: %v", err)
	} else {
		log.Printf("kubernetes client ready")
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	api.RegisterRoutes(r, cfg, client)

	log.Printf("starting %s on %s (env=%s)", cfg.AppName, cfg.Addr, cfg.Env)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
