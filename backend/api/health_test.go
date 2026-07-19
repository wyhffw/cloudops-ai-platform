package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wyhffw/cloudops-ai-platform/backend/pkg/config"
)

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	RegisterRoutes(r, config.Config{
		AppName: "cloudops-backend",
		Env:     "test",
		Version: "0.2.0",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}
}

func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		AppName:   "cloudops-backend",
		Env:       "test",
		Version:   "0.2.0",
		JWTSecret: "test-secret",
		AdminUser: "admin",
		AdminPass: "admin123",
		JWTExpire: time.Hour,
	}

	r := gin.New()
	RegisterRoutes(r, cfg, nil)

	body := bytes.NewBufferString(`{"username":"admin","password":"admin123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["token"] == nil || resp["token"] == "" {
		t.Fatalf("expected token in response")
	}
}

func TestNamespacesRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, config.Config{JWTSecret: "x"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
