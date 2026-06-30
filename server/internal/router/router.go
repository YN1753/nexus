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

func RouterInit(authHandler handler.AuthHandler) *gin.Engine {
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
	}
	return r
}
