package service

import (
	"fmt"
	"nexus/internal/repository"
	pkgssh "nexus/pkg/ssh"
	"nexus/pkg/utils"
)

type SSHService struct {
	NodeRepo repository.NodeRepository
}

func NewSSHService(nodeRepo repository.NodeRepository) SSHService {
	return SSHService{
		NodeRepo: nodeRepo,
	}
}

func (s *SSHService) Connect(nodeID uint) (*pkgssh.Client, error) {
	node, err := s.NodeRepo.FindNodeById(nodeID)
	if err != nil {
		return nil, err
	}

	if node.AuthType == 2 {
		node.PrivateKey, err = utils.Decrypt(node.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("解密私钥失败: %w", err)
		}
	} else {
		node.Password, err = utils.Decrypt(node.Password)
		if err != nil {
			return nil, fmt.Errorf("解密密码失败: %w", err)
		}
	}

	client, err := pkgssh.NewClient(*node)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *SSHService) RunCommand(nodeID uint, cmd string) (string, error) {
	client, err := s.Connect(nodeID)
	if err != nil {
		return "", err
	}
	defer client.Close()

	output, err := client.RunCommand(cmd)
	if err != nil {
		return output, err
	}

	return output, nil
}
