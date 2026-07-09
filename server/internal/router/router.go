package router

import (
	"nexus/internal/handler"
	"nexus/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	c         *gin.Engine
	JwtSecret string
	JwtExpire uint
}

func RouterInit(authHandler handler.AuthHandler, nodeHandler handler.NodeHandler, sshHandler handler.SSHHandler) *gin.Engine {
	r := gin.Default()
	public := r.Group("api/v1")
	{
		public.POST("login", authHandler.Login)
		public.POST("register", authHandler.Register)
	}
	private := r.Group("api/v1")
	{
		private.Use(middleware.JWTAuth(authHandler.JwtSecret))
		private.GET("info", authHandler.GetInfo)

		//node相关的路由
		private.POST("node", nodeHandler.SaveNode)
		private.POST("getNode", nodeHandler.GetNodeById)
		private.GET("nodes", nodeHandler.GetNodes)
		private.POST("updateNode", nodeHandler.UpdateNode)

		// ssh相关的路由
		private.GET("ssh/run", sshHandler.RunCommand)
	}
	return r
}
