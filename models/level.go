package models

import "gitee.com/muzipp/Distribution/pkg/logging"

type Level struct {
	Model
	Name string
}

func ListLevels(maps interface{}, fields string) (levels []Level, flag bool) {

	err := Db.Table("level").
		Where(maps).
		Select(fields).
		Scan(&levels).Error
	if err != nil {
		flag = true
		logging.Info("等级列表错误", err) //记录错误日志
		return
	}
	return
}