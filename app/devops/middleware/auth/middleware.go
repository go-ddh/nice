package auth

import (
	"fmt"
	"github.com/go-ddh/nice/framework/gin"
)

// AuthMiddleware 登录中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件日志信息")
		c.Next()
	}
}
