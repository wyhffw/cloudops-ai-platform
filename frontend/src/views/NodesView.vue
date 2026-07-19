<template>
  <div class="page">
    <div class="toolbar">
      <el-button type="primary" :loading="loading" @click="load">刷新</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe border height="calc(100vh - 180px)">
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.ready ? 'success' : 'danger'" size="small">
            {{ row.ready ? "Ready" : "NotReady" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="roles" label="角色" min-width="140">
        <template #default="{ row }">{{ (row.roles || []).join(", ") }}</template>
      </el-table-column>
      <el-table-column prop="version" label="版本" width="120" />
      <el-table-column prop="os" label="系统" min-width="180" show-overflow-tooltip />
      <el-table-column prop="age" label="Age" width="140" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { fetchNodes } from "@/api/cloudops";

const loading = ref(false);
const items = ref<any[]>([]);

async function load() {
  loading.value = true;
  try {
    const data = await fetchNodes();
    items.value = data.items || [];
  } catch {
    ElMessage.error("加载节点失败");
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
.page {
  background: var(--ops-panel);
  border: 1px solid var(--ops-border);
  border-radius: 12px;
  padding: 16px;
}

.toolbar {
  margin-bottom: 12px;
}
</style>
