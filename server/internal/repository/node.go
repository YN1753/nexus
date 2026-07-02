package repository

import (
	"errors"
	"nexus/internal/model"

	"gorm.io/gorm"
)

type NodeRepository struct {
	DB *gorm.DB
}

func NewNodeRepository(db *gorm.DB) NodeRepository {
	return NodeRepository{
		DB: db,
	}
}

func (repo *NodeRepository) SaveNode(node *model.Node) error {
	err := repo.DB.Create(node).Error
	if err != nil {
		return errors.New("创建节点失败" + err.Error())
	}
	return nil
}
func (repo *NodeRepository) UpdateNode(node *model.Node) error {
	err := repo.DB.Updates(node).Error
	if err != nil {
		return errors.New("更新节点失败" + err.Error())
	}
	return nil
}
func (repo *NodeRepository) DeleteNodeById(id uint) error {
	return repo.DB.Delete(model.Node{}, id).Error
}

func (repo *NodeRepository) FindNodeById(id uint) (*model.Node, error) {
	var node model.Node
	err := repo.DB.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (repo *NodeRepository) GetNodes() ([]model.Node, error) {
	var nodes []model.Node
	err := repo.DB.Find(&nodes).Error
	if err != nil {
		return nil, errors.New("获取nodes失败" + err.Error())
	}
	return nodes, nil
}
