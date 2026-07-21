<script setup lang="ts">
import {
  ArrowLeft,
  Connection,
  Edit,
  FolderOpened,
  Monitor,
  Refresh,
  VideoPlay,
} from '@element-plus/icons-vue'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import NodeFormDialog from '@/components/NodeFormDialog.vue'
import { getNode, getNodeOverview } from '@/services/node'
import type { NodeInfo, NodeOverview } from '@/types/api'

interface MetricCard {
  label: string
  value: string
  detail: string
  percent?: number
}

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

const metricCards = computed<MetricCard[]>(() => {
  if (!overview.value) {
    return []
  }

  return [
    {
      label: 'CPU',
      value: `${overview.value.cpuPercent.toFixed(2)}%`,
      detail: '当前 CPU 使用率',
      percent: overview.value.cpuPercent,
    },
    {
      label: 'Memory',
      value: `${overview.value.memoryPercent.toFixed(2)}%`,
      detail: `${overview.value.memoryUsedMB} MB / ${overview.value.memoryTotalMB} MB`,
      percent: overview.value.memoryPercent,
    },
    {
      label: 'Disk',
      value: `${overview.value.diskPercent.toFixed(2)}%`,
      detail: `${overview.value.diskUsed} / ${overview.value.diskTotal}`,
      percent: overview.value.diskPercent,
    },
    {
      label: 'Load',
      value: overview.value.loadAverage || '-',
      detail: '1m / 5m / 15m',
    },
  ]
})

async function loadNodeDetail() {
  loading.value = true
  try {
    const detail = await getNode({ id: props.nodeId })
    node.value = detail
    await loadOverview()
  } finally {
    loading.value = false
  }
}

async function loadOverview() {
  overviewLoading.value = true
  try {
    overview.value = await getNodeOverview({ id: props.nodeId })
  } finally {
    overviewLoading.value = false
  }
}

function progressStatus(percent?: number) {
  if (typeof percent !== 'number') return undefined
  if (percent >= 80) return 'exception'
  if (percent >= 60) return 'warning'
  return 'success'
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
          <p class="muted">{{ node?.Description || '节点连接信息与系统概览' }}</p>
        </div>

        <div class="title-actions">
          <el-button :icon="Refresh" circle :loading="loading || overviewLoading" @click="loadNodeDetail" />
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
        <div class="summary-kicker">主机信息</div>
        <div class="summary-title">{{ overview?.hostname || '-' }}</div>
        <div class="summary-sub muted">{{ overview?.os || '等待采集系统信息' }}</div>
      </article>

      <article class="summary-panel">
        <div class="summary-kicker">最近采集</div>
        <div class="summary-title">{{ overview?.lastReportAt || '等待采集' }}</div>
        <div class="summary-sub muted">本页数据由后端通过 SSH 实时采集</div>
      </article>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="概览" name="overview">
        <div v-loading="overviewLoading" class="overview-layout">
          <div class="split-layout" v-if="node">
            <section class="content-block">
              <div class="block-header">
                <div>
                  <h2>连接配置</h2>
                  <p class="muted">节点资产与连接参数。</p>
                </div>
                <el-icon class="block-icon"><Monitor /></el-icon>
              </div>

              <div class="system-grid compact-grid">
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
                  <span class="muted">UUID</span>
                  <strong class="mono">{{ node.UUID }}</strong>
                </div>
              </div>
            </section>

            <section class="content-block" v-if="overview">
              <div class="block-header">
                <div>
                  <h2>系统信息</h2>
                  <p class="muted">实时读取远程 Linux 主机状态。</p>
                </div>
                <el-icon class="block-icon"><Connection /></el-icon>
              </div>

              <div class="system-grid compact-grid">
                <div class="system-item">
                  <span class="muted">主机名</span>
                  <strong>{{ overview.hostname || '-' }}</strong>
                </div>
                <div class="system-item">
                  <span class="muted">操作系统</span>
                  <strong>{{ overview.os || '-' }}</strong>
                </div>
                <div class="system-item">
                  <span class="muted">内核</span>
                  <strong>{{ overview.kernel || '-' }}</strong>
                </div>
                <div class="system-item">
                  <span class="muted">架构</span>
                  <strong>{{ overview.architecture || '-' }}</strong>
                </div>
                <div class="system-item">
                  <span class="muted">运行时长</span>
                  <strong>{{ overview.uptime || '-' }}</strong>
                </div>
                <div class="system-item">
                  <span class="muted">平均负载</span>
                  <strong>{{ overview.loadAverage || '-' }}</strong>
                </div>
              </div>
            </section>
          </div>

          <section class="metric-grid" v-if="metricCards.length > 0">
            <article v-for="metric in metricCards" :key="metric.label" class="metric-panel">
              <div class="metric-topline">
                <span class="muted">{{ metric.label }}</span>
                <strong>{{ metric.value }}</strong>
              </div>
              <p class="muted metric-detail">{{ metric.detail }}</p>
              <el-progress
                v-if="typeof metric.percent === 'number'"
                :percentage="metric.percent"
                :status="progressStatus(metric.percent)"
                :show-text="false"
                :stroke-width="8"
              />
            </article>
          </section>

          <section class="content-block" v-if="overview">
            <div class="block-header">
              <div>
                <h2>挂载点</h2>
                <p class="muted">来自远程节点的真实磁盘挂载信息。</p>
              </div>
              <el-icon class="block-icon"><FolderOpened /></el-icon>
            </div>

            <div v-if="overview.mounts.length > 0" class="mount-list">
              <article v-for="mount in overview.mounts" :key="mount.name" class="list-row">
                <div class="list-copy">
                  <strong class="mono">{{ mount.name }}</strong>
                  <span class="muted">{{ mount.fileSystem }} · {{ mount.used }} / {{ mount.size }}</span>
                </div>
                <div class="list-side">
                  <strong>{{ mount.usage }}</strong>
                  <el-progress :percentage="mount.percent" :show-text="false" :stroke-width="8" />
                </div>
              </article>
            </div>
            <el-empty v-else description="暂无挂载点信息" />
          </section>
        </div>
      </el-tab-pane>

      <el-tab-pane label="终端" name="terminal">
        <section class="content-block placeholder-block">
          <el-icon class="placeholder-icon"><Monitor /></el-icon>
          <h2>终端入口</h2>
          <p class="muted">继续使用你现在已经打通的 WebSSH 终端。</p>
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

.compact-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
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

.mount-list {
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
  .split-layout,
  .compact-grid {
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
