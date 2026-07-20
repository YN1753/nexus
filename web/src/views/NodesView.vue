<script setup lang="ts">
import {
  Connection,
  Document,
  Edit,
  Monitor,
  Plus,
  Refresh,
  VideoPlay,
} from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import NodeFormDialog from '@/components/NodeFormDialog.vue'
import { getNodes } from '@/services/node'
import type { NodeInfo } from '@/types/api'

const router = useRouter()
const loading = ref(false)
const nodes = ref<NodeInfo[]>([])
const createVisible = ref(false)
const editVisible = ref(false)
const editingNode = ref<NodeInfo | null>(null)

const onlineCount = computed(() => nodes.value.filter((node) => node.Status === 1).length)
const passwordCount = computed(() => nodes.value.filter((node) => node.AuthType === 1).length)

async function loadNodes() {
  loading.value = true
  try {
    nodes.value = await getNodes()
  } finally {
    loading.value = false
  }
}

function openTerminal(node: NodeInfo) {
  router.push({ name: 'terminal', params: { nodeId: node.ID } })
}

function openDetail(node: NodeInfo) {
  router.push({ name: 'node-detail', params: { nodeId: node.ID } })
}

function openCreateDialog() {
  createVisible.value = true
}

function openEditDialog(node: NodeInfo) {
  editingNode.value = { ...node }
  editVisible.value = true
}

async function handleCreated() {
  await loadNodes()
}

async function handleUpdated() {
  editVisible.value = false
  await loadNodes()
}

onMounted(loadNodes)
</script>

<template>
  <section class="nodes-page">
    <div class="header-row">
      <div>
        <h1>节点</h1>
        <p class="muted">把你现在已经写好的节点接口都收进一个工作台里。</p>
      </div>

      <div class="header-actions">
        <el-button :icon="Plus" type="primary" @click="openCreateDialog">新增节点</el-button>
        <el-tooltip content="刷新节点" placement="bottom">
          <el-button :icon="Refresh" circle :loading="loading" @click="loadNodes" />
        </el-tooltip>
      </div>
    </div>

    <div class="stat-grid">
      <article class="stat-card">
        <span class="muted">节点总数</span>
        <strong>{{ nodes.length }}</strong>
      </article>
      <article class="stat-card">
        <span class="muted">标记在线</span>
        <strong>{{ onlineCount }}</strong>
      </article>
      <article class="stat-card">
        <span class="muted">密码认证</span>
        <strong>{{ passwordCount }}</strong>
      </article>
    </div>

    <el-empty v-if="!loading && nodes.length === 0" description="暂无节点">
      <el-button :icon="Plus" type="primary" @click="openCreateDialog">创建第一台节点</el-button>
    </el-empty>

    <el-table v-else v-loading="loading" :data="nodes" row-key="ID" class="node-table">
      <el-table-column label="节点" min-width="220">
        <template #default="{ row }">
          <div class="node-name">
            <el-icon><Monitor /></el-icon>
            <div>
              <strong>{{ row.Name }}</strong>
              <span class="muted">{{ row.Description || '无备注' }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="地址" min-width="220">
        <template #default="{ row }">
          <span class="mono">{{ row.SshUser }}@{{ row.Host }}:{{ row.Port }}</span>
        </template>
      </el-table-column>

      <el-table-column label="认证" width="110">
        <template #default="{ row }">
          <el-tag :type="row.AuthType === 2 ? 'warning' : 'success'" effect="plain">
            {{ row.AuthType === 2 ? '密钥' : '密码' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="row.Status === 1 ? 'success' : 'info'" effect="plain">
            <el-icon><Connection /></el-icon>
            {{ row.Status === 1 ? '在线' : '未知' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="220" align="right">
        <template #default="{ row }">
          <div class="action-group">
            <el-tooltip content="查看详情" placement="top">
              <el-button :icon="Document" circle @click="openDetail(row)" />
            </el-tooltip>
            <el-tooltip content="编辑节点" placement="top">
              <el-button :icon="Edit" circle @click="openEditDialog(row)" />
            </el-tooltip>
            <el-tooltip content="打开终端" placement="top">
              <el-button :icon="VideoPlay" type="primary" circle @click="openTerminal(row)" />
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <NodeFormDialog v-model="createVisible" mode="create" @success="handleCreated" />

    <NodeFormDialog
      v-model="editVisible"
      mode="edit"
      :initial-data="editingNode"
      @success="handleUpdated"
    />
  </section>
</template>

<style scoped>
.nodes-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.header-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.header-row h1 {
  margin: 0 0 6px;
  font-size: 26px;
}

.header-row p {
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 18px 20px;
  border: 1px solid #dbe2ea;
  border-radius: 8px;
  background: #ffffff;
}

.stat-card strong {
  font-size: 28px;
  line-height: 1;
}

.node-table {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
}

.node-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.node-name .el-icon {
  color: #0f766e;
  font-size: 20px;
}

.node-name div {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

.node-name strong,
.node-name span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.el-tag {
  gap: 4px;
}

.action-group {
  display: inline-flex;
  gap: 8px;
}

@media (max-width: 900px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .header-row {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
