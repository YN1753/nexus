<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

import { saveNode, updateNode } from '@/services/node'
import type { NodeInfo, SaveNodePayload, UpdateNodePayload } from '@/types/api'

interface NodeFormState {
  name: string
  description: string
  host: string
  port: number
  sshUser: string
  authType: number
  password: string
  privateKey: string
  status: number
}

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    mode: 'create' | 'edit'
    initialData?: NodeInfo | null
  }>(),
  {
    initialData: null,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive<NodeFormState>({
  name: '',
  description: '',
  host: '',
  port: 22,
  sshUser: 'root',
  authType: 1,
  password: '',
  privateKey: '',
  status: 0,
})

const dialogTitle = computed(() => (props.mode === 'create' ? '新增节点' : '编辑节点'))
const isPasswordMode = computed(() => form.authType === 1)

const rules: FormRules<NodeFormState> = {
  name: [{ required: true, message: '请输入节点名称', trigger: 'blur' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入 SSH 端口', trigger: 'blur' }],
  sshUser: [{ required: true, message: '请输入 SSH 用户名', trigger: 'blur' }],
  password: [
    {
      validator: (_, value, callback) => {
        if (isPasswordMode.value && !value) {
          callback(new Error('请输入 SSH 密码'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  privateKey: [
    {
      validator: (_, value, callback) => {
        if (!isPasswordMode.value && !value) {
          callback(new Error('请输入私钥内容'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
}

function syncForm() {
  const data = props.initialData
  form.name = data?.Name || ''
  form.description = data?.Description || ''
  form.host = data?.Host || ''
  form.port = data?.Port || 22
  form.sshUser = data?.SshUser || 'root'
  form.authType = data?.AuthType || 1
  form.password = data?.Password || ''
  form.privateKey = data?.PrivateKey || ''
  form.status = data?.Status || 0
}

function closeDialog() {
  emit('update:modelValue', false)
}

async function submit() {
  if (!formRef.value) {
    return
  }

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  submitting.value = true
  try {
    if (props.mode === 'create') {
      const payload: SaveNodePayload = {
        name: form.name,
        description: form.description,
        host: form.host,
        port: Number(form.port),
        ssh_user: form.sshUser,
        auth_type: form.authType,
        password: form.authType === 1 ? form.password : '',
        privateKey: form.authType === 2 ? form.privateKey : '',
      }

      await saveNode(payload)
      ElMessage.success('节点已创建')
    } else {
      if (!props.initialData) {
        ElMessage.error('缺少节点数据，无法更新')
        return
      }

      const payload: UpdateNodePayload = {
        ID: props.initialData.ID,
        Name: form.name,
        Description: form.description,
        UUID: props.initialData.UUID,
        Host: form.host,
        Port: Number(form.port),
        SshUser: form.sshUser,
        AuthType: form.authType,
        Password: form.authType === 1 ? form.password : '',
        PrivateKey: form.authType === 2 ? form.privateKey : '',
        Status: form.status,
      }

      await updateNode(payload)
      ElMessage.success('节点已更新')
    }

    emit('success')
    closeDialog()
  } finally {
    submitting.value = false
  }
}

watch(
  () => [props.modelValue, props.initialData, props.mode],
  ([visible]) => {
    if (visible) {
      syncForm()
      formRef.value?.clearValidate()
    }
  },
  { immediate: true },
)
</script>

<template>
  <el-dialog
    :model-value="modelValue"
    :title="dialogTitle"
    width="min(680px, calc(100vw - 32px))"
    destroy-on-close
    @close="closeDialog"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
      <div class="form-grid">
        <el-form-item label="节点名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：tx-prod" />
        </el-form-item>

        <el-form-item label="SSH 用户" prop="sshUser">
          <el-input v-model="form.sshUser" placeholder="root / ubuntu" />
        </el-form-item>

        <el-form-item class="full-span" label="备注">
          <el-input v-model="form.description" placeholder="这台机器的用途" />
        </el-form-item>

        <el-form-item label="主机地址" prop="host">
          <el-input v-model="form.host" placeholder="IP 或域名" />
        </el-form-item>

        <el-form-item label="SSH 端口" prop="port">
          <el-input-number v-model="form.port" :min="1" :max="65535" class="full-width" />
        </el-form-item>

        <el-form-item label="认证方式" prop="authType">
          <el-segmented
            v-model="form.authType"
            :options="[
              { label: '密码', value: 1 },
              { label: '私钥', value: 2 },
            ]"
            class="full-width"
          />
        </el-form-item>

        <el-form-item v-if="isPasswordMode" label="SSH 密码" prop="password">
          <el-input v-model="form.password" show-password type="password" />
        </el-form-item>

        <el-form-item v-else class="full-span" label="私钥内容" prop="privateKey">
          <el-input v-model="form.privateKey" :rows="8" type="textarea" />
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">保存</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 16px;
}

.full-span {
  grid-column: 1 / -1;
}

.full-width {
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 720px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
