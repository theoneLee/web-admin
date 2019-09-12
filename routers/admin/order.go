package admin

import (
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
	"gitee.com/muzipp/Distribution/service/admin/order"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//订单列表
func ListOrders(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	memberId := com.StrTo(c.DefaultQuery("member_id", "0")).MustInt()
	status := com.StrTo(c.DefaultQuery("status", "0")).MustInt()
	number := c.DefaultQuery("number", "")
	remark := c.DefaultQuery("remark", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	orderField := c.DefaultQuery("order_field", "")
	orderSort := com.StrTo(c.DefaultQuery("order_sort", "0")).MustInt()
	data := make(map[string]interface{})

	orderService := order.Order{
		Offset:     util.GetPage(c),
		Limit:      setting.AppSetting.PageSize,
		MemberId:   memberId,
		Status:     status,
		Number:     number,
		Remark:     remark,
		StartTime:  startTime,
		EndTime:    endTime,
		OrderField: orderField,
		OrderSort:  orderSort,
	}

	code := e.ERROR_SQL_FAIL
	orders, err := orderService.ListOrders()
	total, totalError := orderService.CountOrders()

	if err.Code == 0 && totalError.Code == 0 {
		code = e.SUCCESS
		data["lists"] = orders
		data["total"] = total
	}

	appG.Response(http.StatusOK, code, data)
}
