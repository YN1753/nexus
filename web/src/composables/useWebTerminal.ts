import { FitAddon } from '@xterm/addon-fit'
import { Terminal } from '@xterm/xterm'
import { nextTick, onBeforeUnmount, ref } from 'vue'
import { ElMessage } from 'element-plus'

import { WebSSHClient } from '@/services/webssh'

interface UseWebTerminalOptions {
  nodeId: number
  token: string
}

export function useWebTerminal(options: UseWebTerminalOptions) {
  const status = ref<'idle' | 'connecting' | 'connected' | 'closed'>('idle')
  const terminal = new Terminal({
    cursorBlink: true,
    fontFamily: 'SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace',
    fontSize: 13,
    lineHeight: 1.25,
    theme: {
      background: '#101724',
      foreground: '#d6deeb',
      cursor: '#5eead4',
      selectionBackground: '#334155',
    },
  })
  const fitAddon = new FitAddon()
  let client: WebSSHClient | null = null
  let resizeObserver: ResizeObserver | null = null

  terminal.loadAddon(fitAddon)

  async function mount(container: HTMLElement) {
    terminal.open(container)
    await nextTick()
    fitAddon.fit()

    client = new WebSSHClient({
      nodeId: options.nodeId,
      token: options.token,
      size: { cols: terminal.cols, rows: terminal.rows },
      onOpen: () => {
        status.value = 'connected'
      },
      onClose: () => {
        status.value = 'closed'
      },
      onError: () => {
        ElMessage.error('WebSSH 连接异常')
      },
      onMessage: (event) => {
        if (event.type === 'output' && event.data) {
          terminal.write(event.data)
        }
        if (event.type === 'error' && event.error) {
          terminal.writeln(`\r\n${event.error}`)
        }
      },
    })

    terminal.onData((data) => client?.sendInput(data))
    resizeObserver = new ResizeObserver(() => fit())
    resizeObserver.observe(container)

    status.value = 'connecting'
    client.connect()
  }

  function fit() {
    fitAddon.fit()
    client?.resize({ cols: terminal.cols, rows: terminal.rows })
  }

  function dispose() {
    resizeObserver?.disconnect()
    client?.close()
    terminal.dispose()
    client = null
  }

  onBeforeUnmount(dispose)

  return {
    status,
    mount,
    fit,
    dispose,
  }
}
