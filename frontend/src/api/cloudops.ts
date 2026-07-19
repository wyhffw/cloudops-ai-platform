import http from "./http";

export interface LoginResult {
  token: string;
  token_type: string;
  expires_at: string;
  user: { username: string; role: string };
}

export const login = (username: string, password: string) =>
  http.post<LoginResult>("/api/v1/login", { username, password }).then((r) => r.data);

export const fetchNamespaces = () =>
  http.get("/api/v1/namespaces").then((r) => r.data);

export const fetchNodes = () => http.get("/api/v1/nodes").then((r) => r.data);

export const fetchPods = (namespace?: string) =>
  http
    .get("/api/v1/pods", { params: namespace ? { namespace } : undefined })
    .then((r) => r.data);

export const fetchDeployments = (namespace?: string) =>
  http
    .get("/api/v1/deployments", { params: namespace ? { namespace } : undefined })
    .then((r) => r.data);

export const fetchPodLogs = (namespace: string, name: string, tail = 200) =>
  http
    .get(`/api/v1/pods/${namespace}/${name}/logs`, { params: { tail } })
    .then((r) => r.data);

export const restartPod = (namespace: string, name: string) =>
  http.post(`/api/v1/pods/${namespace}/${name}/restart`).then((r) => r.data);

export const scaleDeployment = (namespace: string, name: string, replicas: number) =>
  http
    .post(`/api/v1/deployments/${namespace}/${name}/scale`, { replicas })
    .then((r) => r.data);
