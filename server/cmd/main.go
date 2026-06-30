package main

import (
	"nexus/internal/config"
	"nexus/internal/database"
	"nexus/internal/handler"
	"nexus/internal/repository"
	"nexus/internal/router"
	"nexus/internal/service"
)

func main() {
	Config := config.ConfigInit()
	db := database.MysqlInit(Config.Mysql)
	rdb := database.RedisInit(Config.Redis)
	repo := repository.NewUserRepository(db, rdb)
	userService := service.NewUserService(repo)
	authHandler := handler.NewAuthHandler(userService, Config.Jwt.Secret, Config.Jwt.ExpireHour)
	r := router.RouterInit(authHandler)
	r.Run(":" + Config.Server.Port)
}
