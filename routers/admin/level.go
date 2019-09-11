package admin

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/service/admin/level"
	"github.com/gin-gonic/gin"
	"net/http"
)

//商品列表
func ListLevels(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	data := make([]models.Level, 7)

	code := e.ERROR_SQL_FAIL
	levels, err := level.ListLevels()

	if err.Code == 0 {
		code = e.SUCCESS
		data = levels
	}

	appG.Response(http.StatusOK, code, data)
}
