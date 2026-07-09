export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface LoginPayload {
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
