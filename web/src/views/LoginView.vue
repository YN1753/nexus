<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
})

async function submit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    await auth.login(form)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <section class="login-panel">
      <div class="login-copy">
        <p class="eyebrow">Nexus Workspace</p>
        <h1>服务器工作台</h1>
        <p class="muted">管理节点、打开终端、执行日常运维操作。</p>
      </div>

      <el-form class="login-form" label-position="top" @submit.prevent="submit">
        <el-form-item label="用户名">
          <el-input v-model="form.username" autocomplete="username" size="large" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            autocomplete="current-password"
            show-password
            size="large"
            type="password"
            @keyup.enter="submit"
          />
        </el-form-item>
        <el-button class="submit" type="primary" size="large" :loading="loading" @click="submit">
          登录
        </el-button>
      </el-form>
    </section>
  </div>
</template>

<style scoped>
.login-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  padding: 24px;
  background: #eef3f7;
}

.login-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  width: min(860px, 100%);
  min-height: 420px;
  overflow: hidden;
  border: 1px solid #d9e0e8;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 20px 60px rgb(23 37 84 / 10%);
}

.login-copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 48px;
  background: #172033;
  color: #ffffff;
}

.eyebrow {
  margin: 0 0 12px;
  color: #5eead4;
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
}

h1 {
  margin: 0 0 16px;
  font-size: 34px;
  line-height: 1.15;
}

.login-copy .muted {
  color: #c9d4e3;
  line-height: 1.7;
}

.login-form {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 40px;
}

.submit {
  width: 100%;
  margin-top: 8px;
}

@media (max-width: 720px) {
  .login-panel {
    grid-template-columns: 1fr;
  }

  .login-copy {
    padding: 32px;
  }

  .login-form {
    padding: 32px;
  }
}
</style>
