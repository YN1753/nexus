import { request } from './http'
import type {
  NodeDetailRequest,
  NodeInfo,
  NodeOverview,
  SaveNodePayload,
  UpdateNodePayload,
} from '@/types/api'

export function getNodes() {
  return request<NodeInfo[]>({
    url: '/nodes',
    method: 'GET',
  })
}

export function getNode(payload: NodeDetailRequest) {
  return request<NodeInfo>({
    url: '/getNode',
    method: 'POST',
    data: payload,
  })
}

export function saveNode(payload: SaveNodePayload) {
  return request<null>({
    url: '/node',
    method: 'POST',
    data: payload,
  })
}

export function updateNode(payload: UpdateNodePayload) {
  return request<null>({
    url: '/updateNode',
    method: 'POST',
    data: payload,
  })
}

export function getNodeOverview(payload: NodeDetailRequest) {
  return request<NodeOverview>({
    url: '/getNodeOverview',
    method: 'POST',
    data: payload,
  })
}
