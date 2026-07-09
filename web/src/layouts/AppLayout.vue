<script setup lang="ts">
import { Monitor, SwitchButton } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.replace({ name: 'login' })
}
</script>

<template>
  <div class="page app-layout">
    <header class="topbar">
      <div class="page-shell topbar-inner">
        <RouterLink class="brand" :to="{ name: 'nodes' }">
          <el-icon><Monitor /></el-icon>
          <span>Nexus</span>
        </RouterLink>
        <div class="topbar-actions">
          <span class="muted">{{ auth.user?.Username || 'Developer' }}</span>
          <el-tooltip content="退出登录" placement="bottom">
            <el-button :icon="SwitchButton" circle @click="logout" />
          </el-tooltip>
        </div>
      </div>
    </header>

    <main class="page-shell main-content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
}

.topbar {
  border-bottom: 1px solid #e4e7ed;
  background: #ffffff;
}

.topbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: #18212f;
  font-size: 18px;
  font-weight: 700;
  text-decoration: none;
}

.brand .el-icon {
  color: #0f766e;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.main-content {
  flex: 1;
  padding: 24px 0;
}
</style>
