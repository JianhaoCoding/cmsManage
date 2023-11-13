package middleware

import (
	"cms/helpers"
	"github.com/gin-gonic/gin"
	"strings"
)

// 空接口返回Data
var resData map[string]interface{}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 过滤不需要验证的请求
		if c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		// 获取jwtToken
		token := c.Request.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			helpers.EndRequest(c, 401, resData, "No authorization header provided")
			c.Abort()
			return
		}

		// 去除前缀Bearer
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := helpers.ParseToken(token)
		if err != nil {
			helpers.EndRequest(c, 401, resData, "Invalid token")
			c.Abort()
			return
		}
		c.Set("adminer_id", claims.UserId)

		// @todo...验证当前用户是否拥有该请求权限
		// @todo...当然也需要效验请求是否合法比如ip、来源地址等等

		c.Next()
	}
}
