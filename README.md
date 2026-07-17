# CloudOps AI Platform

Enterprise cloud-native operations platform based on Kubernetes, Go, and AI Agents.

## Structure

- `frontend/` — Vue 3 + Element Plus
- `backend/` — Go (Gin) API
- `k8s/` — Kubernetes manifests (GitOps)
- `.github/workflows/` — GitHub Actions CI/CD

## Backend (quick start)

```bash
cd backend
go test ./...
go run ./cmd/server
```

Then open:

- http://127.0.0.1:8080/healthz
- http://127.0.0.1:8080/api/v1/info

Optional env:

- `ADDR` (default `:8080`)
- `APP_ENV` (`dev` / `prod`)

Build image (needs Docker):

```bash
cd backend
docker build -t cloudops-backend:0.1.0 .
docker run --rm -p 8080:8080 cloudops-backend:0.1.0
```
