<script setup lang="ts">
import { ArrowLeft, RefreshRight } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { useWebTerminal } from '@/composables/useWebTerminal'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  nodeId: number
}>()

const router = useRouter()
const auth = useAuthStore()
const terminalRef = ref<HTMLElement | null>(null)
const terminal = useWebTerminal({ nodeId: props.nodeId, token: auth.token })

const statusText = computed(() => {
  const map = {
    idle: '待连接',
    connecting: '连接中',
    connected: '已连接',
    closed: '已关闭',
  }
  return map[terminal.status.value]
})

onMounted(() => {
  if (terminalRef.value) {
    terminal.mount(terminalRef.value)
  }
})
</script>

<template>
  <section class="terminal-page">
    <div class="header-row">
      <div>
        <el-button :icon="ArrowLeft" link @click="router.push({ name: 'nodes' })">
          返回节点
        </el-button>
        <h1>终端 #{{ nodeId }}</h1>
      </div>
      <div class="terminal-actions">
        <el-tag :type="terminal.status.value === 'connected' ? 'success' : 'info'" effect="plain">
          {{ statusText }}
        </el-tag>
        <el-tooltip content="适配窗口" placement="bottom">
          <el-button :icon="RefreshRight" circle @click="terminal.fit" />
        </el-tooltip>
      </div>
    </div>

    <div class="terminal-frame">
      <div ref="terminalRef" class="terminal-host" />
    </div>
  </section>
</template>

<style scoped>
.terminal-page {
  display: flex;
  min-height: calc(100vh - 112px);
  flex-direction: column;
  gap: 16px;
}

.header-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

h1 {
  margin: 4px 0 0;
  font-size: 24px;
}

.terminal-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.terminal-frame {
  display: flex;
  min-height: 560px;
  flex: 1;
  overflow: hidden;
  border: 1px solid #1f2937;
  border-radius: 8px;
  background: #101724;
}

.terminal-host {
  width: 100%;
  min-width: 0;
  min-height: 0;
  padding: 12px;
}

@media (max-width: 720px) {
  .header-row {
    align-items: flex-start;
    flex-direction: column;
  }

  .terminal-frame {
    min-height: 520px;
  }
}
</style>
