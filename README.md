# CloudOps AI Platform

Enterprise cloud-native operations platform based on Kubernetes, Go, and AI Agents.

## Structure

- `frontend/` — Vue 3 + Element Plus
- `backend/` — Go (Gin) API
- `k8s/` — Kubernetes manifests (GitOps)
- `ansible/` — 两节点 kubeadm 集群安装剧本
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

## CI/CD

Push to `main` triggers GitHub Actions:

1. `go test` / `go build`
2. Build image and push to `ghcr.io/wyhffw/cloudops-backend` (`latest` + short SHA)

Harbor can replace GHCR later; manifests already use a registry image tag.

## Deploy to Kubernetes

```bash
kubectl apply -k k8s/apps/backend
kubectl -n cloudops get pods,svc
kubectl -n cloudops port-forward svc/backend 8080:8080
```

Then open http://127.0.0.1:8080/healthz

> First deploy needs the image on GHCR (after CI succeeds). If the package is private, create an imagePullSecret for GHCR.
