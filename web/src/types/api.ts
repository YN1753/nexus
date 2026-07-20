export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface LoginPayload {
  username: string
  password: string
}

export interface RegisterPayload {
  username: string
  password: string
}

export interface UserInfo {
  ID: number
  Username: string
  CreatedAt?: string
}

export interface NodeInfo {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: null | string
  Name: string
  Description: string
  UUID: string
  Host: string
  Port: number
  SshUser: string
  AuthType: number
  Password?: string
  PrivateKey?: string
  Status: number
}

export interface NodeDetailRequest {
  id: number
}

export interface SaveNodePayload {
  name: string
  description: string
  host: string
  port: number
  ssh_user: string
  auth_type: number
  password: string
  privateKey: string
}

export interface UpdateNodePayload {
  ID: number
  Name: string
  Description: string
  UUID?: string
  Host: string
  Port: number
  SshUser: string
  AuthType: number
  Password: string
  PrivateKey: string
  Status?: number
}

export interface SystemMetric {
  label: string
  value: string
  detail: string
  percent?: number
  tone?: 'default' | 'success' | 'warning' | 'danger'
}

export interface NodeOverview {
  os: string
  kernel: string
  architecture: string
  uptime: string
  loadAverage: string
  lastReportAt: string
  metrics: SystemMetric[]
  mounts: Array<{
    name: string
    usage: string
    detail: string
    percent: number
  }>
  services: Array<{
    name: string
    status: string
    note: string
    tone: 'success' | 'warning' | 'info'
  }>
}
