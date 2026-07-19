<template>
  <div class="page">
    <div class="toolbar">
      <el-select v-model="namespace" clearable filterable placeholder="全部 Namespace" style="width: 220px" @change="load">
        <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
      </el-select>
      <el-button type="primary" :loading="loading" @click="load">刷新</el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe border height="calc(100vh - 180px)">
      <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
      <el-table-column prop="namespace" label="Namespace" width="140" />
      <el-table-column prop="phase" label="状态" width="110" />
      <el-table-column prop="node" label="节点" width="120" />
      <el-table-column prop="restarts" label="重启" width="80" />
      <el-table-column prop="age" label="Age" width="140" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openLogs(row)">日志</el-button>
          <el-button link type="danger" @click="onRestart(row)">重启</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-drawer v-model="logsVisible" size="50%" :title="logsTitle">
      <pre class="logs">{{ logs }}</pre>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { fetchNamespaces, fetchPodLogs, fetchPods, restartPod } from "@/api/cloudops";

const loading = ref(false);
const items = ref<any[]>([]);
const namespaces = ref<string[]>([]);
const namespace = ref<string>("cloudops");
const logsVisible = ref(false);
const logsTitle = ref("");
const logs = ref("");

async function loadNamespaces() {
  const data = await fetchNamespaces();
  namespaces.value = (data.items || []).map((i: any) => i.name);
}

async function load() {
  loading.value = true;
  try {
    const data = await fetchPods(namespace.value || undefined);
    items.value = data.items || [];
  } catch {
    ElMessage.error("加载 Pods 失败");
  } finally {
    loading.value = false;
  }
}

async function openLogs(row: any) {
  logsTitle.value = `${row.namespace}/${row.name}`;
  logsVisible.value = true;
  logs.value = "加载中...";
  try {
    const data = await fetchPodLogs(row.namespace, row.name, 200);
    logs.value = data.logs || "(empty)";
  } catch {
    logs.value = "加载日志失败";
  }
}

async function onRestart(row: any) {
  await ElMessageBox.confirm(`确认重启 Pod ${row.name}？`, "提示", { type: "warning" });
  try {
    await restartPod(row.namespace, row.name);
    ElMessage.success("已触发重启");
    await load();
  } catch {
    ElMessage.error("重启失败");
  }
}

onMounted(async () => {
  try {
    await loadNamespaces();
  } catch {
    /* ignore */
  }
  await load();
});
</script>

<style scoped>
.page {
  background: var(--ops-panel);
  border: 1px solid var(--ops-border);
  border-radius: 12px;
  padding: 16px;
}

.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}

.logs {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, "Courier New", monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #d7e2ef;
}
</style>
