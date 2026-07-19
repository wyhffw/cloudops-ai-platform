<template>
  <el-container class="shell">
    <el-aside width="220px" class="aside">
      <div class="logo">CloudOps</div>
      <el-menu :default-active="route.path" router background-color="#171d25" text-color="#c9d4e3" active-text-color="#2dd4bf">
        <el-menu-item index="/nodes">
          <el-icon><Monitor /></el-icon>
          <span>Nodes</span>
        </el-menu-item>
        <el-menu-item index="/pods">
          <el-icon><Box /></el-icon>
          <span>Pods</span>
        </el-menu-item>
        <el-menu-item index="/deployments">
          <el-icon><CopyDocument /></el-icon>
          <span>Deployments</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="title">{{ title }}</div>
        <div class="right">
          <span class="user">{{ auth.username }}</span>
          <el-button text type="danger" @click="onLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

const title = computed(() => {
  const map: Record<string, string> = {
    "/nodes": "节点 Nodes",
    "/pods": "容器组 Pods",
    "/deployments": "工作负载 Deployments",
  };
  return map[route.path] || "CloudOps";
});

function onLogout() {
  auth.logout();
  router.push("/login");
}
</script>

<style scoped>
.shell {
  min-height: 100vh;
}

.aside {
  border-right: 1px solid var(--ops-border);
  background: var(--ops-panel);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: var(--ops-accent);
  border-bottom: 1px solid var(--ops-border);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--ops-border);
  background: #121820;
}

.title {
  font-size: 18px;
  font-weight: 600;
}

.right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user {
  color: var(--ops-muted);
}

.main {
  background: var(--ops-bg);
}
</style>
