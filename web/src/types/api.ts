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

export interface NodeOverview {
  hostname: string
  os: string
  kernel: string
  architecture: string
  uptime: string
  loadAverage: string
  cpuPercent: number
  memoryTotalMB: number
  memoryUsedMB: number
  memoryPercent: number
  diskTotal: string
  diskUsed: string
  diskPercent: number
  lastReportAt: string
  mounts: Array<{
    name: string
    fileSystem: string
    size: string
    used: string
    usage: string
    percent: number
  }>
}
