package repository

import (
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
	return repo.DB.Create(node).Error
}
func (repo *NodeRepository) UpdateNode(node *model.Node) error {
	return repo.DB.Save(node).Error
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
