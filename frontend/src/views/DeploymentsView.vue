<template>
  <div class="page">
    <div class="toolbar">
      <el-select v-model="namespace" clearable filterable placeholder="全部 Namespace" style="width: 220px" @change="load">
        <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
      </el-select>
      <el-button type="primary" :loading="loading" @click="load">刷新</el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe border height="calc(100vh - 180px)">
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="namespace" label="Namespace" width="140" />
      <el-table-column prop="ready" label="Ready" width="100" />
      <el-table-column prop="replicas" label="副本" width="90" />
      <el-table-column prop="available" label="可用" width="90" />
      <el-table-column prop="age" label="Age" width="140" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="onScale(row)">扩缩容</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { fetchDeployments, fetchNamespaces, scaleDeployment } from "@/api/cloudops";

const loading = ref(false);
const items = ref<any[]>([]);
const namespaces = ref<string[]>([]);
const namespace = ref("cloudops");

async function loadNamespaces() {
  const data = await fetchNamespaces();
  namespaces.value = (data.items || []).map((i: any) => i.name);
}

async function load() {
  loading.value = true;
  try {
    const data = await fetchDeployments(namespace.value || undefined);
    items.value = data.items || [];
  } catch {
    ElMessage.error("加载 Deployments 失败");
  } finally {
    loading.value = false;
  }
}

async function onScale(row: any) {
  const { value } = await ElMessageBox.prompt("输入目标副本数", `扩缩容 ${row.name}`, {
    inputValue: String(row.replicas ?? 1),
    inputPattern: /^\d+$/,
    inputErrorMessage: "请输入非负整数",
  });
  try {
    await scaleDeployment(row.namespace, row.name, Number(value));
    ElMessage.success("已提交扩缩容");
    await load();
  } catch {
    ElMessage.error("扩缩容失败");
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
</style>
