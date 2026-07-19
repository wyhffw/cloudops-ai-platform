package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
	"k8s.io/client-go/kubernetes"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config, client *kubernetes.Clientset) {
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
			"version": cfg.Version,
			"env":     cfg.Env,
		})
	})

	v1 := r.Group("/api/v1")
	registerAuthRoutes(v1, cfg)

	protected := v1.Group("")
	protected.Use(JWTAuth(cfg))
	if client != nil {
		registerK8sRoutes(protected, client)
	} else {
		protected.GET("/namespaces", func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kubernetes client unavailable"})
		})
		protected.GET("/pods", func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kubernetes client unavailable"})
		})
	}
}
