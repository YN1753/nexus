import { request } from './http'
import type { LoginPayload, RegisterPayload, UserInfo } from '@/types/api'

export function login(payload: LoginPayload) {
  return request<string>({
    url: '/login',
    method: 'POST',
    data: payload,
  })
}

export function register(payload: RegisterPayload) {
  return request<number>({
    url: '/register',
    method: 'POST',
    data: payload,
  })
}

export function getUserInfo() {
  return request<UserInfo>({
    url: '/info',
    method: 'GET',
  })
}
