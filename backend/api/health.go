package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.AppName,
			"env":     cfg.Env,
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	r.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    cfg.AppName,
			"version": "0.1.0",
			"env":     cfg.Env,
		})
	})
}
