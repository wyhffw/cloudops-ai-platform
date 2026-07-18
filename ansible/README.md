# Ansible：两节点 Kubernetes 安装

目标主机：

| 主机 | IP |
|------|-----|
| master | 192.168.1.20 |
| worker1 | 192.168.1.21 |

## 推荐：在 master 上当控制节点跑

### 1. 在 master 安装 Ansible

```bash
sudo apt update
sudo apt install -y ansible git
```

### 2. 配置 SSH 免密（master → 自己 + worker1）

```bash
ssh-keygen -t ed25519 -N '' -f ~/.ssh/id_ed25519
ssh-copy-id root@192.168.1.20
ssh-copy-id root@192.168.1.21
```

若安装时创建的是普通用户，改 `inventory/hosts.ini` 里的 `ansible_user`。

### 3. 拿到剧本

任选一种：

```bash
# 从 GitHub 拉
git clone git@github.com:wyhffw/cloudops-ai-platform.git
cd cloudops-ai-platform/ansible

# 或用 U 盘 / scp 把本机 ansible 目录拷到 master
```

### 4. 先测连通

```bash
cd ansible
ansible all -m ping
```

应看到两台 `SUCCESS` / `pong`。

### 5. 一键安装

```bash
ansible-playbook site.yml
```

### 6. 验证

在 master：

```bash
kubectl get nodes
kubectl get pods -A
```

两台都是 `Ready` 即成功。

## 分步执行（可选）

```bash
ansible-playbook site.yml --tags never   # 不适用时忽略
# 或只跑前置：先改 site 为分 playbook；当前用整包 site.yml 即可
```

## 常见问题

1. **ping 失败**：检查 `ansible_user`、SSH 密钥、IP  
2. **拉包慢/失败**：已用阿里云镜像；确认虚拟机能上网  
3. **worker join 失败**：重新在 master 执行  
   `kubeadm token create --print-join-command`  
   再单独对 worker 跑 join  
4. **重复执行**：已 init/join 的节点会跳过关键步骤（idempotent 尽量保证）
