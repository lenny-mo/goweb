package middleware

import (
	"go_web_app/controllers"
	"go_web_app/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWT() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 获取authorization header
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// 没有携带token
			controllers.ReturnResponse(c, http.StatusUnauthorized, controllers.NeedAuthCode)
			c.Abort() // 终止中间件函数后续的调用
			return
		}
		// 2. 按空格分割
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// token格式不对
			controllers.ReturnResponse(c, http.StatusUnauthorized, controllers.InvalidTokenCode)
			c.Abort()
			return
		}

		// 4. 解析令牌
		claims, err := pkg.ParseToken(parts[1])
		if err != nil {
			// 解析token失败
			controllers.ReturnResponse(c, http.StatusUnauthorized, controllers.InvalidTokenCode)
			c.Abort() // 终止中间件函数后续的调用
			return
		}
		// 5. 获取claims中的username和userID
		// 6. 将username和userID绑定到请求上下文中，后续的处理函数就可以用c.Get("username")来获取当前请求的用户信息
		c.Set(controllers.ContextUsernameKey, claims.Username)
		c.Set(controllers.ContextUserIDKey, claims.UserID)
		c.Set(controllers.ExpireTimeKey, *claims.ExpiresAt)
		c.Next() // 后续的处理函数可以用c.Get("username")来获取当前请求的用户信息
	}
}
