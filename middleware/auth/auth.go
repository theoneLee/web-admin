package auth

import (
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/gredis"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
返回一个HandlerFunc类型值
*/
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		/**
		初始化返回码200
		*/
		code = e.SUCCESS
		token := c.GetHeader("token")

		/**
		判断token是否为空
		*/
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			//从redis中根据token获取value，判断token是否失效
			redisUser, err := gredis.Get("user_token" + token)
			if redisUser == nil || err != nil { //token验证失败
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}

		/**
		验证不通过的情况，直接返回错误信息
		*/
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
