package e

/**
声明MsgFlags为map（映射），key为int类型，value为string类型
*/
var MsgFlags = map[int]string{
	SUCCESS:                        "操作成功",
	ERROR:                          "操作失败",
	INVALID_PARAMS:                 "参数不合法",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "token失效，请重新登录",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "鉴权失败",
	ERROR_AUTH_TOKEN:               "token生成失败",
	ERROR_AUTH:                     "用户名或密码不合法",
	ERROR_REDIS:                    "Redis设置失败",
	ERROR_SQL_FAIL:                 "操作数据库失败",
}

/**
传来的code在MsgFlags中不存在的情况，会返回Token错误
*/
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
