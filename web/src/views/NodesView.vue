<script setup lang="ts">
import { Connection, Monitor, Refresh, VideoPlay } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getNodes } from '@/services/node'
import type { NodeInfo } from '@/types/api'

const router = useRouter()
const loading = ref(false)
const nodes = ref<NodeInfo[]>([])

const onlineCount = computed(() => nodes.value.filter((node) => node.Status === 1).length)

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

onMounted(loadNodes)
</script>

<template>
  <section class="nodes-page">
    <div class="header-row">
      <div>
        <h1>节点</h1>
        <p class="muted">{{ nodes.length }} 台服务器，{{ onlineCount }} 台标记在线。</p>
      </div>
      <el-tooltip content="刷新节点" placement="bottom">
        <el-button :icon="Refresh" circle :loading="loading" @click="loadNodes" />
      </el-tooltip>
    </div>

    <el-empty v-if="!loading && nodes.length === 0" description="暂无节点" />

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

      <el-table-column label="操作" width="120" align="right">
        <template #default="{ row }">
          <el-tooltip content="打开终端" placement="top">
            <el-button :icon="VideoPlay" type="primary" circle @click="openTerminal(row)" />
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>
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
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

h1 {
  margin: 0 0 6px;
  font-size: 26px;
}

.header-row p {
  margin: 0;
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
</style>
