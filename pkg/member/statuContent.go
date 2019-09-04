package member

var SexFlags = map[int]string{
	SEX_MALE:"男",
	SEX_FEMALE:"女",
}

var StatusFlags = map[int]string{
	STATUS_NOT_ACTIVE:"未激活",
	STATUS_EXPRIRED:"已过期",
	STATUS_DELETE:"已删除",
	STATUS_DISABLE:"已禁用",
	STATUS_NORMAL:"正常",
}


func GetSex(code int) string {
	msg, ok := SexFlags[code]
	if ok {
		return msg
	}

	return SexFlags[code]
}


func GetStatus(code int) string {
	msg, ok := StatusFlags[code]
	if ok {
		return msg
	}

	return StatusFlags[code]
}
