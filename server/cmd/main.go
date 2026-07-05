package main

import (
	"nexus/internal/config"
	"nexus/internal/database"
	"nexus/internal/handler"
	"nexus/internal/repository"
	"nexus/internal/router"
	"nexus/internal/service"
	"nexus/pkg/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	utils.AesKey = []byte(os.Getenv("AES_KEY"))
	Config := config.ConfigInit()
	db := database.MysqlInit(Config.Mysql)
	rdb := database.RedisInit(Config.Redis)
	repo := repository.NewUserRepository(db, rdb)
	userService := service.NewUserService(repo)
	nodeRepo := repository.NewNodeRepository(db)
	nodeService := service.NewNodeService(nodeRepo)
	sshService := service.NewSSHService(nodeRepo)
	authHandler := handler.NewAuthHandler(userService, Config.Jwt.Secret, Config.Jwt.ExpireHour)
	nodeHandler := handler.NewNodeHandler(nodeService)
	sshHandler := handler.NewSSHHandler(sshService)
	r := router.RouterInit(authHandler, nodeHandler, sshHandler)
	r.Run(":" + Config.Server.Port)
}
