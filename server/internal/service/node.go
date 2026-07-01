package service

import (
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
	return n.Repo.UpdateNode(node)
}
