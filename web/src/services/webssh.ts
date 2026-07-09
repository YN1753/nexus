export type TerminalEventType = 'ready' | 'output' | 'error' | 'close'

export interface TerminalEvent {
  type: TerminalEventType
  data?: string
  error?: string
}

export interface TerminalSize {
  cols: number
  rows: number
}

export interface WebSSHClientOptions {
  nodeId: number
  token: string
  size?: TerminalSize
  onMessage?: (event: TerminalEvent) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (event: Event) => void
}

export class WebSSHClient {
  private socket: WebSocket | null = null

  constructor(private readonly options: WebSSHClientOptions) {}

  connect() {
    this.socket = new WebSocket(this.buildURL())

    this.socket.onopen = () => this.options.onOpen?.()
    this.socket.onclose = () => this.options.onClose?.()
    this.socket.onerror = (event) => this.options.onError?.(event)
    this.socket.onmessage = (event) => this.handleMessage(event.data)
  }

  sendInput(data: string) {
    this.send({ type: 'input', data })
  }

  resize(size: TerminalSize) {
    this.send({ type: 'resize', rows: size.rows, cols: size.cols })
  }

  close() {
    this.send({ type: 'close' })
    this.socket?.close()
    this.socket = null
  }

  private send(payload: Record<string, unknown>) {
    if (this.socket?.readyState !== WebSocket.OPEN) {
      return
    }
    this.socket.send(JSON.stringify(payload))
  }

  private handleMessage(raw: string) {
    try {
      const event = JSON.parse(raw) as TerminalEvent
      this.options.onMessage?.(event)
    } catch {
      this.options.onMessage?.({ type: 'output', data: raw })
    }
  }

  private buildURL() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const params = new URLSearchParams({
      node_id: String(this.options.nodeId),
      token: this.options.token,
    })

    if (this.options.size) {
      params.set('cols', String(this.options.size.cols))
      params.set('rows', String(this.options.size.rows))
    }

    return `${protocol}//${window.location.host}/api/v1/ssh/run?${params.toString()}`
  }
}
