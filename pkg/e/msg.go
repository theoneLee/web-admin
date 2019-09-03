package e

/**
声明MsgFlags为map（映射），key为int类型，value为string类型
*/
var MsgFlags = map[int]string {
	SUCCESS : "ok",
	ERROR : "fail",
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