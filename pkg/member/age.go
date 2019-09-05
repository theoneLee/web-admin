package member

import (
	"fmt"
	"time"
)

func GetTimeFromStrDate(date string) (year, month, day int) {
	const shortForm = "2006/01/02"
	d, err := time.Parse(shortForm, date)
	if err != nil {
		fmt.Println("出生日期解析错误！")
		return 0, 0, 0
	}
	year = d.Year()
	month = int(d.Month())
	day = d.Day()
	return
}

func GetAge(year int) (age int) {
	if year <= 0 {
		age = -1
	}
	nowYear := time.Now().Year()
	age = nowYear - year
	return
}