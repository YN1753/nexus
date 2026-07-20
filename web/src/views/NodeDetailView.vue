<script setup lang="ts">
import {
  ArrowLeft,
  Connection,
  Edit,
  FolderOpened,
  Monitor,
  Refresh,
  TrendCharts,
  VideoPlay,
} from '@element-plus/icons-vue'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import NodeFormDialog from '@/components/NodeFormDialog.vue'
import { getMockNodeOverview, getNode } from '@/services/node'
import type { NodeInfo, NodeOverview, SystemMetric } from '@/types/api'

const props = defineProps<{
  nodeId: number
}>()

const router = useRouter()
const loading = ref(false)
const overviewLoading = ref(false)
const activeTab = ref('overview')
const editVisible = ref(false)
const node = ref<NodeInfo | null>(null)
const overview = ref<NodeOverview | null>(null)

const statusText = computed(() => (node.value?.Status === 1 ? '在线' : '未知'))
const authText = computed(() => (node.value?.AuthType === 2 ? '密钥认证' : '密码认证'))

async function loadNodeDetail() {
  loading.value = true
  try {
    const detail = await getNode({ id: props.nodeId })
    node.value = detail
    await loadOverview(detail)
  } finally {
    loading.value = false
  }
}

async function loadOverview(detail: NodeInfo) {
  overviewLoading.value = true
  try {
    overview.value = await getMockNodeOverview(detail)
  } finally {
    overviewLoading.value = false
  }
}

function metricStatus(metric: SystemMetric) {
  if (metric.tone === 'success') return 'success'
  if (metric.tone === 'warning') return 'warning'
  if (metric.tone === 'danger') return 'exception'
  return undefined
}

function goTerminal() {
  router.push({ name: 'terminal', params: { nodeId: props.nodeId } })
}

async function handleUpdated() {
  editVisible.value = false
  await loadNodeDetail()
}

watch(
  () => props.nodeId,
  () => {
    loadNodeDetail()
  },
)

onMounted(loadNodeDetail)
</script>

<template>
  <section class="detail-page" v-loading="loading">
    <div class="detail-header">
      <el-button :icon="ArrowLeft" link @click="router.push({ name: 'nodes' })">返回节点</el-button>

      <div class="title-row">
        <div class="title-copy">
          <h1>{{ node?.Name || `节点 #${nodeId}` }}</h1>
          <p class="muted">{{ node?.Description || '开发节点概览与运行状态' }}</p>
        </div>

        <div class="title-actions">
          <el-button :icon="Refresh" circle :loading="loading" @click="loadNodeDetail" />
          <el-button :icon="Edit" @click="editVisible = true">编辑节点</el-button>
          <el-button :icon="VideoPlay" type="primary" @click="goTerminal">打开终端</el-button>
        </div>
      </div>
    </div>

    <div v-if="node" class="summary-grid">
      <article class="summary-panel summary-primary">
        <div class="summary-kicker">节点身份</div>
        <div class="summary-title">{{ node.SshUser }}@{{ node.Host }}:{{ node.Port }}</div>
        <div class="summary-meta">
          <el-tag :type="node.Status === 1 ? 'success' : 'info'" effect="plain">
            <el-icon><Connection /></el-icon>
            {{ statusText }}
          </el-tag>
          <el-tag :type="node.AuthType === 2 ? 'warning' : 'success'" effect="plain">
            {{ authText }}
          </el-tag>
        </div>
      </article>

      <article class="summary-panel">
        <div class="summary-kicker">UUID</div>
        <div class="summary-title mono small">{{ node.UUID }}</div>
        <div class="summary-sub muted">用于资产标识与后续能力扩展</div>
      </article>

      <article class="summary-panel">
        <div class="summary-kicker">最近同步</div>
        <div class="summary-title">{{ overview?.lastReportAt || '等待采集' }}</div>
        <div class="summary-sub muted">当前仍为前端模拟数据</div>
      </article>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="概览" name="overview">
        <div v-loading="overviewLoading" class="overview-layout">
          <section class="content-block" v-if="node">
            <div class="block-header">
              <div>
                <h2>连接配置</h2>
                <p class="muted">这里展示你后端真实返回的节点详情。</p>
              </div>
              <el-icon class="block-icon"><Monitor /></el-icon>
            </div>

            <div class="system-grid">
              <div class="system-item">
                <span class="muted">名称</span>
                <strong>{{ node.Name }}</strong>
              </div>
              <div class="system-item">
                <span class="muted">主机地址</span>
                <strong>{{ node.Host }}</strong>
              </div>
              <div class="system-item">
                <span class="muted">端口</span>
                <strong>{{ node.Port }}</strong>
              </div>
              <div class="system-item">
                <span class="muted">SSH 用户</span>
                <strong>{{ node.SshUser }}</strong>
              </div>
              <div class="system-item">
                <span class="muted">认证方式</span>
                <strong>{{ authText }}</strong>
              </div>
              <div class="system-item">
                <span class="muted">创建时间</span>
                <strong>{{ node.CreatedAt }}</strong>
              </div>
            </div>
          </section>

          <section class="metric-grid" v-if="overview">
            <article v-for="metric in overview.metrics" :key="metric.label" class="metric-panel">
              <div class="metric-topline">
                <span class="muted">{{ metric.label }}</span>
                <strong>{{ metric.value }}</strong>
              </div>
              <p class="muted metric-detail">{{ metric.detail }}</p>
              <el-progress
                v-if="typeof metric.percent === 'number'"
                :percentage="metric.percent"
                :status="metricStatus(metric)"
                :show-text="false"
                :stroke-width="8"
              />
            </article>
          </section>

          <div class="split-layout" v-if="overview">
            <section class="content-block">
              <div class="block-header">
                <div>
                  <h2>挂载点</h2>
                  <p class="muted">现在先用模拟数据占位，之后可以接真实采集。</p>
                </div>
                <el-icon class="block-icon"><FolderOpened /></el-icon>
              </div>

              <div class="mount-list">
                <article v-for="mount in overview.mounts" :key="mount.name" class="list-row">
                  <div class="list-copy">
                    <strong class="mono">{{ mount.name }}</strong>
                    <span class="muted">{{ mount.detail }}</span>
                  </div>
                  <div class="list-side">
                    <strong>{{ mount.usage }}</strong>
                    <el-progress :percentage="mount.percent" :show-text="false" :stroke-width="8" />
                  </div>
                </article>
              </div>
            </section>

            <section class="content-block">
              <div class="block-header">
                <div>
                  <h2>服务状态</h2>
                  <p class="muted">也还是前端模拟，让布局先稳定下来。</p>
                </div>
                <el-icon class="block-icon"><TrendCharts /></el-icon>
              </div>

              <div class="service-list">
                <article v-for="service in overview.services" :key="service.name" class="list-row service-row">
                  <div class="list-copy">
                    <strong class="mono">{{ service.name }}</strong>
                    <span class="muted">{{ service.note }}</span>
                  </div>
                  <el-tag :type="service.tone" effect="plain">{{ service.status }}</el-tag>
                </article>
              </div>
            </section>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="终端" name="terminal">
        <section class="content-block placeholder-block">
          <el-icon class="placeholder-icon"><Monitor /></el-icon>
          <h2>终端入口</h2>
          <p class="muted">你现在已经把 WebSSH 跑通了，这里保留节点内的主入口。</p>
          <el-button :icon="VideoPlay" type="primary" @click="goTerminal">进入终端</el-button>
        </section>
      </el-tab-pane>
    </el-tabs>

    <NodeFormDialog v-model="editVisible" mode="edit" :initial-data="node" @success="handleUpdated" />
  </section>
</template>

<style scoped>
.detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-header {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.title-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.title-copy h1 {
  margin: 0;
  font-size: 28px;
}

.title-copy p {
  margin: 8px 0 0;
}

.title-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.summary-grid {
  display: grid;
  grid-template-columns: 1.4fr 1fr 1fr;
  gap: 16px;
}

.summary-panel,
.content-block,
.metric-panel {
  border: 1px solid #dbe2ea;
  border-radius: 8px;
  background: #ffffff;
}

.summary-panel {
  padding: 20px;
}

.summary-primary {
  background: linear-gradient(135deg, #172033 0%, #20304a 100%);
  border-color: #172033;
  color: #ffffff;
}

.summary-primary .muted,
.summary-primary .summary-kicker,
.summary-primary .summary-sub {
  color: #c4d0e3;
}

.summary-kicker {
  margin-bottom: 10px;
  color: #667085;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
}

.summary-title {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.4;
}

.summary-title.small {
  font-size: 14px;
  line-height: 1.6;
}

.summary-sub {
  margin-top: 8px;
}

.summary-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}

.overview-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-block {
  padding: 20px;
}

.block-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.block-header h2,
.placeholder-block h2 {
  margin: 0;
  font-size: 20px;
}

.block-header p,
.placeholder-block p {
  margin: 6px 0 0;
}

.block-icon,
.placeholder-icon {
  color: #0f766e;
  font-size: 20px;
}

.system-grid,
.metric-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.system-item,
.metric-panel {
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #ffffff;
}

.system-item {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 8px;
}

.system-item strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.metric-topline {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 10px;
}

.metric-topline strong {
  font-size: 24px;
  line-height: 1;
}

.metric-detail {
  min-height: 22px;
  margin: 10px 0 14px;
}

.split-layout {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.mount-list,
.service-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
}

.list-copy,
.list-side {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
}

.list-side {
  width: min(180px, 100%);
  align-items: flex-end;
}

.service-row {
  align-items: flex-start;
}

.placeholder-block {
  display: flex;
  min-height: 320px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  text-align: center;
}

@media (max-width: 960px) {
  .summary-grid,
  .system-grid,
  .metric-grid,
  .split-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .title-row,
  .title-actions,
  .list-row {
    align-items: flex-start;
    flex-direction: column;
  }

  .list-side {
    width: 100%;
    align-items: flex-start;
  }
}
</style>
