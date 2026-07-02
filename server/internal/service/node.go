package service

import (
	"errors"
	"nexus/internal/model"
	"nexus/internal/repository"
	"nexus/pkg/utils"
)

type NodeService struct {
	Repo repository.NodeRepository
}

func NewNodeService(repo repository.NodeRepository) NodeService {
	return NodeService{
		Repo: repo,
	}
}

func (n *NodeService) SaveNode(node *model.Node) error {
	node.UUID = utils.GetUUID()
	var err error
	node.Password, err = utils.Encrypt(node.Password)
	if err != nil {
		return err
	}
	node.PrivateKey, err = utils.Encrypt(node.PrivateKey)
	if err != nil {
		return err
	}
	return n.Repo.SaveNode(node)
}

func (n *NodeService) GetNode(id uint) (*model.Node, error) {
	node, err := n.Repo.FindNodeById(id)
	if err != nil {
		return nil, err
	}
	node.Password, err = utils.Decrypt(node.Password)
	if err != nil {
		return nil, err
	}
	node.PrivateKey, err = utils.Decrypt(node.PrivateKey)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (n *NodeService) UpdateNode(node *model.Node) error {
	var err error
	node.Password, err = utils.Encrypt(node.Password)
	if err != nil {
		return err
	}
	node.PrivateKey, err = utils.Encrypt(node.PrivateKey)
	if err != nil {
		return err
	}
	return n.Repo.UpdateNode(node)
}

func (n *NodeService) GetNodes() (*[]model.Node, error) {
	nodes, err := n.Repo.GetNodes()
	if err != nil {
		return nil, err
	}
	for index, _ := range nodes {
		nodes[index].Password, err = utils.Decrypt(nodes[index].Password)
		if err != nil {
			return nil, errors.New("信息解码失败")
		}
		nodes[index].PrivateKey, err = utils.Decrypt(nodes[index].PrivateKey)
		if err != nil {
			return nil, errors.New("信息解码失败")
		}
	}
	return &nodes, nil
}
