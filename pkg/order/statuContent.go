package order


var StatusFlags = map[int]string{
	ORDER_STATUS_EXPRIRED:   "无效已关闭",
	ORDER_STATUS_DELETE:     "待审核",
	ORDER_STATUS_DISABLE:    "待发货",
	ORDER_STATUS_NORMAL:     "已发货",
}

func GetStatus(code int) string {
	msg, ok := StatusFlags[code]
	if ok {
		return msg
	}

	return StatusFlags[code]
}
