package handler

import (
	"nexus/internal/model"
	"nexus/internal/request"
	"nexus/internal/service"
	"nexus/pkg/response"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	NodeService service.NodeService
}

func NewNodeHandler(nodeService service.NodeService) NodeHandler {
	return NodeHandler{
		NodeService: nodeService,
	}
}

func (n *NodeHandler) SaveNode(c *gin.Context) {
	var req request.SaveNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "获取json失败")
		return
	}
	node := model.Node{
		Name:        req.Name,
		Description: req.Description,
		Password:    req.Password,
		Host:        req.Host,
		Port:        req.Port,
		SshUser:     req.SshUser,
		AuthType:    req.AuthType,
		PrivateKey:  req.PrivateKey,
	}
	err := n.NodeService.SaveNode(&node)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
	return
}
func (n *NodeHandler) GetNodeById(c *gin.Context) {
	var req request.GetNodeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, 400, "获取json错误")
		return
	}
	node, err := n.NodeService.GetNode(req.ID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, &node)
	return
}
