package ssh

import (
	"errors"
	"fmt"
	"nexus/internal/model"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type Client struct {
	client *gossh.Client
}

func NewClient(node model.Node) (*Client, error) {
	if node.AuthType != 1 && node.AuthType != 2 {
		return nil, fmt.Errorf("不支持的认证类型: %d", node.AuthType)
	}

	var authMethods []gossh.AuthMethod

	if node.AuthType == 2 {
		signer, err := gossh.ParsePrivateKey([]byte(node.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
		authMethods = append(authMethods, gossh.PublicKeys(signer))
	} else {
		authMethods = append(authMethods, gossh.Password(node.Password))
	}

	config := &gossh.ClientConfig{
		User:            node.SshUser,
		Auth:            authMethods,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", node.Host, node.Port)
	client, err := gossh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) RunCommand(cmd string) (string, error) {
	if c == nil || c.client == nil {
		return "", errors.New("ssh客户端未初始化")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建ssh会话失败: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("执行命令失败: %w", err)
	}

	return string(output), nil
}

func (c *Client) NewSession() (*gossh.Session, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("ssh客户端未初始化")
	}
	return c.client.NewSession()
}

func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}
