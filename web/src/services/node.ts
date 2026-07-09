import { request } from './http'
import type { NodeInfo } from '@/types/api'

export function getNodes() {
  return request<NodeInfo[]>({
    url: '/nodes',
    method: 'GET',
  })
}
