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

export async function getMockNodeOverview(node: NodeInfo): Promise<NodeOverview> {
  const seed = node.ID || 1
  const cpu = 18 + ((seed * 13) % 42)
  const memory = 34 + ((seed * 9) % 38)
  const disk = 26 + ((seed * 7) % 48)
  const network = 12 + ((seed * 11) % 58)

  await new Promise((resolve) => window.setTimeout(resolve, 220))

  return {
    os: 'Ubuntu 24.04.4 LTS',
    kernel: 'Linux 6.8.0-101-generic',
    architecture: 'x86_64',
    uptime: `${6 + seed} days ${3 + (seed % 6)} hours`,
    loadAverage: `0.${seed} / 0.${seed + 2} / 0.${seed + 4}`,
    lastReportAt: '2026-07-20 20:30:00',
    metrics: [
      {
        label: 'CPU',
        value: `${cpu}%`,
        detail: '8 vCPU · current load',
        percent: cpu,
        tone: cpu > 75 ? 'danger' : cpu > 55 ? 'warning' : 'success',
      },
      {
        label: 'Memory',
        value: `${memory}%`,
        detail: '8.2 GB / 16 GB',
        percent: memory,
        tone: memory > 75 ? 'danger' : memory > 55 ? 'warning' : 'success',
      },
      {
        label: 'Disk',
        value: `${disk}%`,
        detail: '146 GB / 512 GB',
        percent: disk,
        tone: disk > 80 ? 'danger' : disk > 60 ? 'warning' : 'success',
      },
      {
        label: 'Network',
        value: `${network} MB/s`,
        detail: 'eth0 outbound peak',
        percent: Math.min(network, 100),
        tone: 'default',
      },
    ],
    mounts: [
      {
        name: '/',
        usage: `${disk}% used`,
        detail: 'ext4 · 512 GB total',
        percent: disk,
      },
      {
        name: '/var/lib/docker',
        usage: `${Math.min(disk + 8, 92)}% used`,
        detail: 'overlay2 · container layers',
        percent: Math.min(disk + 8, 92),
      },
      {
        name: '/data',
        usage: `${Math.max(disk - 12, 10)}% used`,
        detail: 'xfs · project workspace',
        percent: Math.max(disk - 12, 10),
      },
    ],
    services: [
      {
        name: 'docker.service',
        status: 'running',
        note: 'containers healthy',
        tone: 'success',
      },
      {
        name: 'redis.service',
        status: 'running',
        note: 'port 6379 listening',
        tone: 'success',
      },
      {
        name: 'nexus-agent-preview',
        status: 'pending',
        note: 'reserved for future worker',
        tone: 'info',
      },
    ],
  }
}
