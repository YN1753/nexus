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
	SSHService  service.SSHService
}

func NewNodeHandler(nodeService service.NodeService, sshService service.SSHService) NodeHandler {
	return NodeHandler{
		NodeService: nodeService,
		SSHService:  sshService,
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

func (n *NodeHandler) GetNodes(c *gin.Context) {
	nodes, err := n.NodeService.GetNodes()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nodes)
	return
}

func (n *NodeHandler) GetNodeOverview(c *gin.Context) {
	var req request.GetNodeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, 400, "获取json错误")
		return
	}
	overview, err := n.SSHService.GetNodeOverview(req.ID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, overview)
	return
}

func (n *NodeHandler) UpdateNode(c *gin.Context) {
	node := model.Node{}
	if err := c.ShouldBindJSON(&node); err != nil {
		response.Fail(c, 400, "json获取失败")
		return
	}
	err := n.NodeService.UpdateNode(&node)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
	return
}
