package middleware

import (
	"nexus/pkg/response"
	"nexus/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			} else {
				response.Fail(c, 401, "Token 格式错误，必须是 Bearer 模式")
				c.Abort()
				return
			}
		} else {
			token = c.Query("token")
		}

		if token == "" {
			response.Fail(c, 401, "请求未携带Token，访问被拒绝")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token, secret)
		if err != nil {
			response.Fail(c, 401, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
