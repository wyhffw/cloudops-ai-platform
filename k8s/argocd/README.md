# Argo CD（装在家里的 K8s 集群）

## 1. 安装

在 master 上：

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 等待就绪
kubectl -n argocd get pods
kubectl -n argocd wait --for=condition=Available deploy/argocd-server --timeout=300s
```

## 2. 获取初始密码

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

用户名：`admin`

## 3. 访问 UI（本机先用端口转发）

```bash
kubectl -n argocd port-forward svc/argocd-server 8081:443 --address 0.0.0.0
```

浏览器打开：`https://192.168.30.20:8081`（接受自签证书）  
登录：`admin` + 上面密码。

## 4. 注册应用

```bash
kubectl apply -f k8s/argocd/application-backend.yaml
```

或在 UI：New App → 仓库 `https://github.com/wyhffw/cloudops-ai-platform.git`，路径 `k8s/apps/backend`，目标 namespace `cloudops`。

## 5. 验证

```bash
kubectl -n argocd get application
kubectl -n cloudops get pods,svc,ingress
```

GitHub 上改 `k8s/apps/backend` 并 push 后，Argo CD 会自动同步（约 3 分钟内，或点 Sync）。
