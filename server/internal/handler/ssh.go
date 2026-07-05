package handler

import (
	"nexus/internal/request"
	"nexus/internal/service"
	"nexus/pkg/response"

	"github.com/gin-gonic/gin"
)

type SSHHandler struct {
	SSHService service.SSHService
}

func NewSSHHandler(sshService service.SSHService) SSHHandler {
	return SSHHandler{
		SSHService: sshService,
	}
}

func (h *SSHHandler) RunCommand(c *gin.Context) {
	var req request.RunCommandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "获取json失败")
		return
	}

	output, err := h.SSHService.RunCommand(req.NodeID, req.Command)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, output)
}
