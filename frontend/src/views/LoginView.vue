<template>
  <div class="login-page">
    <div class="panel">
      <div class="brand">CloudOps</div>
      <p class="sub">Kubernetes 运维控制台</p>
      <el-form :model="form" @submit.prevent="onSubmit">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            show-password
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-button type="primary" size="large" class="submit" :loading="loading" @click="onSubmit">
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const loading = ref(false);
const form = reactive({ username: "admin", password: "admin123" });

async function onSubmit() {
  loading.value = true;
  try {
    await auth.login(form.username, form.password);
    ElMessage.success("登录成功");
    const redirect = (route.query.redirect as string) || "/nodes";
    await router.replace(redirect);
  } catch (e: unknown) {
    ElMessage.error("登录失败，请检查账号密码或后端服务");
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at 20% 20%, rgba(45, 212, 191, 0.12), transparent 35%),
    radial-gradient(circle at 80% 0%, rgba(56, 189, 248, 0.1), transparent 30%),
    linear-gradient(160deg, #0b1117, #121a22 50%, #0f1419);
}

.panel {
  width: min(400px, 92vw);
  padding: 36px 32px 28px;
  border: 1px solid var(--ops-border);
  border-radius: 16px;
  background: rgba(23, 29, 37, 0.92);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
}

.brand {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.04em;
  color: var(--ops-accent);
}

.sub {
  margin: 8px 0 28px;
  color: var(--ops-muted);
}

.submit {
  width: 100%;
  margin-top: 8px;
}
</style>
