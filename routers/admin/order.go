package admin

import (
	"fmt"
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
	"gitee.com/muzipp/Distribution/routers/common"
	"gitee.com/muzipp/Distribution/service/admin/order"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
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

//订单详情
func DetailOrder(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	//验证有没有错误
	if valid.HasErrors() {
		//记录验证错误日志
		app.MarkErrors(valid.Errors)
		//请求返回
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	code := e.ERROR_SQL_FAIL
	var data map[string]interface{}

	//获取商品详情
	orderService := order.Order{
		Id: id,
	}

	orderRst, err := orderService.DetailOrder()
	if err.Code == 0 {
		code = e.SUCCESS
		data = orderRst
	}

	appG.Response(http.StatusOK, code, data)

}

func AddOrder(c *gin.Context) {
	appG := app.Gin{C: c}                               //实例化响应对象
	memberName := c.DefaultPostForm("member_name", "0") //订单所属人ID
	remark := c.DefaultPostForm("remark", "")           //备注
	billNumber := c.DefaultPostForm("bill_number", "")  //银行流水号
	goodsInfo := c.DefaultPostForm("goods_info", "")    //下单的商品信息

	valid := validation.Validation{}
	valid.Required(memberName, "member_name").Message("订单所属人不能为空")
	valid.Required(billNumber, "bill_number").Message("银行流水号不能为空")
	valid.Required(goodsInfo, "goods_info").Message("商品信息不能为空")

	//设置返回数据
	data := make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		//图片上传
		orderService := order.Order{
			Remark:      remark,
			MemberName:  memberName,
			BillNumber:  billNumber,
			GoodsInfo:   goodsInfo,
			RecommendId: common.SelfUser.Id,
		}
		err := orderService.AddOrder()

		if err.Code == 0 {
			code = e.SUCCESS
		} else {
			code = err.Code
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info(fmt.Sprintf("%s,%s", "err key is "+err.Key, "err Message is "+err.Message))
		}
	}

	appG.Response(http.StatusOK, code, data)
}

//会员列表
func OrderStatusChange(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.PostForm("id")).MustInt()
	status := com.StrTo(c.PostForm("status")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Range(status, -3, 1, "status").Message("状态只允许-3或-1")

	code := e.INVALID_PARAMS
	if status != 1 && status != -3 {
		goto End
	}
	if !valid.HasErrors() {
		orderService := order.Order{
			Id:     id,
			Status: status,
		}
		code = e.ERROR_SQL_FAIL
		err := orderService.StatusChange()
		if err.Code == 0 {
			code = e.SUCCESS
		} else {
			code = err.Code
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

End:
	appG.Response(http.StatusOK, code, nil)

}
