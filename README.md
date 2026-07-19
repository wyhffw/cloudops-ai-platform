# CloudOps AI Platform

Enterprise cloud-native operations platform based on Kubernetes, Go, and AI Agents.

## Structure

- `frontend/` — Vue 3 + Element Plus
- `backend/` — Go (Gin) API
- `k8s/` — Kubernetes manifests (GitOps)
- `ansible/` — 两节点 kubeadm 集群安装剧本
- `.github/workflows/` — GitHub Actions CI/CD

## Frontend (quick start)

```bash
# 终端 1：后端
cd backend && go run ./cmd/server

# 终端 2：前端
cd frontend
npm install
npm run dev
```

打开 http://127.0.0.1:5173 ，默认账号 `admin` / `admin123`。
Vite 已把 `/api` 代理到后端 `:8080`。

```bash
cd backend
go test ./...
go run ./cmd/server
```

Then open:

- http://127.0.0.1:8080/healthz
- http://127.0.0.1:8080/api/v1/info

Auth + cluster APIs (need kubeconfig or in-cluster SA):

```bash
# login
curl -s http://127.0.0.1:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'

# list namespaces / pods / nodes / deployments (replace TOKEN)
curl -s http://127.0.0.1:8080/api/v1/namespaces -H "Authorization: Bearer TOKEN"
curl -s 'http://127.0.0.1:8080/api/v1/pods?namespace=cloudops' -H "Authorization: Bearer TOKEN"
curl -s http://127.0.0.1:8080/api/v1/nodes -H "Authorization: Bearer TOKEN"
curl -s 'http://127.0.0.1:8080/api/v1/deployments?namespace=cloudops' -H "Authorization: Bearer TOKEN"

# pod logs / restart / scale
curl -s 'http://127.0.0.1:8080/api/v1/pods/cloudops/POD_NAME/logs?tail=100' -H "Authorization: Bearer TOKEN"
curl -s -X POST 'http://127.0.0.1:8080/api/v1/pods/cloudops/POD_NAME/restart' -H "Authorization: Bearer TOKEN"
curl -s -X POST 'http://127.0.0.1:8080/api/v1/deployments/cloudops/backend/scale' \
  -H "Authorization: Bearer TOKEN" -H 'Content-Type: application/json' \
  -d '{"replicas":2}'
```

Optional env:

- `ADDR` (default `:8080`)
- `APP_ENV` (`dev` / `prod`)
- `JWT_SECRET` / `ADMIN_USER` / `ADMIN_PASSWORD`
- `KUBECONFIG` (local only; in-cluster uses ServiceAccount)

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
