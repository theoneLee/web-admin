package goods

var StatusFlags = map[int]string{
	STATUS_STATUS_NORMAL:  "正常",
	STATUS_GOODS_OUT:      "缺货",
	STATUS_GOODS_OBTAINED: "下架",
}

func GetStatus(code int) string {
	msg, ok := StatusFlags[code]
	if ok {
		return msg
	}

	return StatusFlags[code]
}
