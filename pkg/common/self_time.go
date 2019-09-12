package common

import "time"

func GetTime(selfTime string) (formatTime int64) {
	const shortForm = "2006/01/02"
	stamp, _ := time.ParseInLocation(shortForm, selfTime, time.Local)
	formatTime = stamp.Unix()
	return
}
