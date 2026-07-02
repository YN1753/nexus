package repository

import (
	"errors"
	"nexus/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
	return UserRepository{
		DB:  db,
		Rdb: rdb,
	}
}

func (u *UserRepository) CreateUser(user model.User) error {
	err := u.DB.Create(&user).Error
	if err != nil {
		return errors.New("创建用户失败" + err.Error())
	}
	return nil
}

func (u *UserRepository) FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, errors.New("查找用户失败（username）" + err.Error())
	}
	return user, err
}

func (u *UserRepository) FindUserById(id uint) (model.User, error) {
	var user model.User
	err := u.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, errors.New("查找用户失败（id）" + err.Error())
	}
	return user, err
}

func (u *UserRepository) UpdateUserInfo(user model.User) error {
	err := u.DB.Updates(&user).Error
	if err != nil {
		return errors.New("更新用户信息失败" + err.Error())
	}
	return nil
}
