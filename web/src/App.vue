<script setup lang="ts">
import { onMounted } from 'vue'

import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

onMounted(() => {
  if (auth.token && !auth.user) {
    auth.loadUser().catch(() => undefined)
  }
})
</script>

<template>
  <RouterView />
</template>
