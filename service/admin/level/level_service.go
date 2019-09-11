package level

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
)

func ListLevels() (levels []models.Level, err e.SelfError) {
	fields := "id,name"
	levels, levelErr := models.ListLevels(getMaps(), fields)
	if levelErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return

}

//封装搜索条件
func getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_at"] = 0
	return maps
}
