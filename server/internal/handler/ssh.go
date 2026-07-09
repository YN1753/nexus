package handler

import (
	"fmt"
	"io"
	"net/http"
	"nexus/internal/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	gossh "golang.org/x/crypto/ssh"
)

type SSHHandler struct {
	SSHService service.SSHService
}

func NewSSHHandler(sshService service.SSHService) SSHHandler {
	return SSHHandler{
		SSHService: sshService,
	}
}

var sshUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type terminalMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
}

type terminalResult struct {
	Type  string `json:"type"`
	Data  string `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type wsWriter struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (w *wsWriter) writeJSON(v any) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.conn.WriteJSON(v)
}

func (w *wsWriter) Write(p []byte) (int, error) {
	if err := w.writeJSON(terminalResult{Type: "output", Data: string(p)}); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (h *SSHHandler) RunCommand(c *gin.Context) {
	conn, err := sshUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	writer := &wsWriter{conn: conn}
	if err := writer.writeJSON(terminalResult{Type: "ready"}); err != nil {
		return
	}

	nodeID, err := parseUintQuery(c, "node_id", 0)
	if err != nil || nodeID == 0 {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: "node_id is required"})
		return
	}

	cols, err := parseUintQuery(c, "cols", 80)
	if err != nil || cols == 0 {
		cols = 80
	}
	rows, err := parseUintQuery(c, "rows", 24)
	if err != nil || rows == 0 {
		rows = 24
	}

	client, err := h.SSHService.Connect(uint(nodeID))
	if err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}
	defer session.Close()

	modes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", int(rows), int(cols), modes); err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}

	if err := session.Shell(); err != nil {
		_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
		return
	}

	done := make(chan struct{})
	var once sync.Once
	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	go func() {
		_, copyErr := io.Copy(writer, stdout)
		if copyErr != nil {
			_ = writer.writeJSON(terminalResult{Type: "error", Error: copyErr.Error()})
		}
		closeDone()
	}()

	go func() {
		_, copyErr := io.Copy(writer, stderr)
		if copyErr != nil {
			_ = writer.writeJSON(terminalResult{Type: "error", Error: copyErr.Error()})
		}
		closeDone()
	}()

	for {
		select {
		case <-done:
			return
		default:
		}

		var msg terminalMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}

		switch msg.Type {
		case "input":
			if _, err := io.WriteString(stdin, msg.Data); err != nil {
				_ = writer.writeJSON(terminalResult{Type: "error", Error: err.Error()})
				return
			}
		case "resize":
			if msg.Rows > 0 && msg.Cols > 0 {
				_ = session.WindowChange(int(msg.Rows), int(msg.Cols))
			}
		case "close":
			return
		default:
			_ = writer.writeJSON(terminalResult{Type: "error", Error: fmt.Sprintf("不支持的消息类型: %s", msg.Type)})
		}
	}
}

func parseUintQuery(c *gin.Context, key string, defaultValue uint64) (uint64, error) {
	raw := c.Query(key)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s 参数格式错误: %w", key, err)
	}

	return value, nil
}
