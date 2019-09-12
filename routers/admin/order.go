package admin

import (
	"fmt"
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
	"gitee.com/muzipp/Distribution/service/admin/goods"
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
	appG := app.Gin{C: c} //实例化响应对象
	name := c.DefaultPostForm("name", "")
	stock := com.StrTo(c.DefaultPostForm("stock", "1")).MustInt()
	status := com.StrTo(c.DefaultPostForm("status", "1")).MustInt()
	integral := com.StrTo(c.DefaultPostForm("integral", "0")).MustInt() //积分
	specification := c.DefaultPostForm("specification", "")             //规格
	price := com.StrTo(c.DefaultPostForm("price", "1")).MustFloat64()
	remark := c.DefaultPostForm("remark", "")
	description := c.DefaultPostForm("description", "") //描述
	image := c.DefaultPostForm("image", "")             //图片地址

	valid := validation.Validation{}
	valid.Required(stock, "stock").Message("库存不能为空")
	valid.Required(name, "name").Message("名称不能为空")
	valid.Required(price, "price").Message("单价不能为空")
	valid.Required(specification, "specification").Message("规格不能为空")
	valid.Required(image, "image").Message("图片不能为空")

	//设置返回数据
	data := make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		//图片上传
		goodsService := goods.Goods{
			Name:          name,
			Stock:         stock,
			Remark:        remark,
			Price:         price,
			Status:        status,
			Specification: specification,
			Integral:      integral,
			Description:   description,
			Image:         image,
		}
		err := goodsService.AddGoods()

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
