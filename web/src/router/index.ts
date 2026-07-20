import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/AppLayout.vue'),
      children: [
        {
          path: '',
          name: 'nodes',
          component: () => import('@/views/NodesView.vue'),
        },
        {
          path: 'nodes/:nodeId',
          name: 'node-detail',
          component: () => import('@/views/NodeDetailView.vue'),
          props: (route) => ({ nodeId: Number(route.params.nodeId) }),
        },
        {
          path: 'terminal/:nodeId',
          name: 'terminal',
          component: () => import('@/views/TerminalView.vue'),
          props: (route) => ({ nodeId: Number(route.params.nodeId) }),
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (!to.meta.public && !auth.isLogin) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.name === 'login' && auth.isLogin) {
    return { name: 'nodes' }
  }

  return true
})

export default router
