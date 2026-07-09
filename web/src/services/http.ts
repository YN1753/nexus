import axios, { AxiosError } from 'axios'
import { ElMessage } from 'element-plus'

import { getToken } from './storage'
import type { ApiResponse } from '@/types/api'

export const http = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse<unknown>
    if (body.code !== 200) {
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return response
  },
  (error: AxiosError<ApiResponse<null>>) => {
    const message = error.response?.data?.message || error.message || '网络请求失败'
    ElMessage.error(message)
    return Promise.reject(new Error(message))
  },
)

export async function request<T>(config: Parameters<typeof http.request>[0]) {
  const response = await http.request<ApiResponse<T>>(config)
  return response.data.data
}
