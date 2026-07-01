package model

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);not null;comment:节点别名"`
	Description string `gorm:"type:varchar(256);comment:备注描述"`
	UUID        string `gorm:"type:varchar(64);not null;comment:UUID"`
	Host        string `gorm:"type:varchar(128);not null;index:idx_host_port,unique;comment:IP或域名"`
	Port        int    `gorm:"type:int;default:22;index:idx_host_port,unique;comment:SSH端口"`

	SshUser    string `gorm:"type:varchar(64);default:'root';not null;comment:SSH登录用户名"`
	AuthType   int    `gorm:"type:tinyint;default:1;comment:认证类型: 1-密码, 2-密钥"`
	Password   string `gorm:"type:text;comment:AES加密后的SSH密码"`
	PrivateKey string `gorm:"type:text;comment:AES加密后的SSH私钥内容"`

	Status int `gorm:"type:tinyint;default:0;comment:在线状态: 0-未知/离线, 1-在线"`
}

func (Node) TableName() string {
	return "nodes"
}
