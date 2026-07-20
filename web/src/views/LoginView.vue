<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
})

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
})

async function submitLogin() {
  if (!loginForm.username || !loginForm.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    await auth.login(loginForm)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  } finally {
    loading.value = false
  }
}

async function submitRegister() {
  if (!registerForm.username || !registerForm.password) {
    ElMessage.warning('请先填写完整注册信息')
    return
  }

  if (registerForm.password !== registerForm.confirmPassword) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    await auth.register({
      username: registerForm.username,
      password: registerForm.password,
    })

    loginForm.username = registerForm.username
    loginForm.password = registerForm.password
    registerForm.confirmPassword = ''
    activeTab.value = 'login'
    ElMessage.success('注册成功，请登录')
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
        <p class="muted">登录、管理节点、进入终端，把你当前后端已经支持的能力先完整跑起来。</p>
      </div>

      <div class="login-form-wrap">
        <el-tabs v-model="activeTab" class="auth-tabs" stretch>
          <el-tab-pane label="登录" name="login">
            <el-form class="login-form" label-position="top" @submit.prevent="submitLogin">
              <el-form-item label="用户名">
                <el-input v-model="loginForm.username" autocomplete="username" size="large" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input
                  v-model="loginForm.password"
                  autocomplete="current-password"
                  show-password
                  size="large"
                  type="password"
                  @keyup.enter="submitLogin"
                />
              </el-form-item>
              <el-button
                class="submit"
                type="primary"
                size="large"
                :loading="loading"
                @click="submitLogin"
              >
                登录
              </el-button>
            </el-form>
          </el-tab-pane>

          <el-tab-pane label="注册" name="register">
            <el-form class="login-form" label-position="top" @submit.prevent="submitRegister">
              <el-form-item label="用户名">
                <el-input v-model="registerForm.username" autocomplete="username" size="large" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input
                  v-model="registerForm.password"
                  autocomplete="new-password"
                  show-password
                  size="large"
                  type="password"
                />
              </el-form-item>
              <el-form-item label="确认密码">
                <el-input
                  v-model="registerForm.confirmPassword"
                  autocomplete="new-password"
                  show-password
                  size="large"
                  type="password"
                  @keyup.enter="submitRegister"
                />
              </el-form-item>
              <el-button
                class="submit"
                size="large"
                :loading="loading"
                @click="submitRegister"
              >
                注册
              </el-button>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
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
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 420px);
  width: min(980px, 100%);
  min-height: 480px;
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

.login-form-wrap {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 28px 32px;
}

.auth-tabs {
  width: 100%;
}

.login-form {
  padding-top: 12px;
}

.submit {
  width: 100%;
  margin-top: 8px;
}

@media (max-width: 780px) {
  .login-panel {
    grid-template-columns: 1fr;
  }

  .login-copy,
  .login-form-wrap {
    padding: 32px;
  }
}
</style>
