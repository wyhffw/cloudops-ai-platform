package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/auth"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func registerAuthRoutes(r *gin.RouterGroup, cfg config.Config) {
	r.POST("/login", func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
			return
		}
		if req.Username != cfg.AdminUser || req.Password != cfg.AdminPass {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, exp, err := auth.IssueToken(cfg.JWTSecret, req.Username, "admin", cfg.JWTExpire)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"token_type": "Bearer",
			"expires_at": exp.UTC().Format(time.RFC3339),
			"user": gin.H{
				"username": req.Username,
				"role":     "admin",
			},
		})
	})
}
