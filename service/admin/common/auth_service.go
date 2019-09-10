package common

import (
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/gredis"
	"strconv"
)

//退出登录功能
func Logout(token string, id int) (err e.SelfError) {

	_, tokenErr := gredis.Delete("user_token" + token)
	_, userErr := gredis.Delete("user_id" + strconv.Itoa(id))

	if tokenErr != nil || userErr != nil {
		err.Code = e.ERROR_REDIS
	}

	return
}

