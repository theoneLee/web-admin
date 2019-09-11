package admin

import (
	"fmt"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
	"gitee.com/muzipp/Distribution/service/admin/goods"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//添加商品
func AddGoods(c *gin.Context) {
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

//商品列表
func ListGoods(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	data := make(map[string]interface{})

	goodsService := goods.Goods{
		Offset: util.GetPage(c),
		Limit:  setting.AppSetting.PageSize,
	}

	code := e.ERROR_SQL_FAIL
	members, err := goodsService.ListGoods()
	total, totalError := goodsService.CountGoods()

	if err.Code == 0 && totalError.Code == 0 {
		code = e.SUCCESS
		data["lists"] = members
		data["total"] = total
	}

	appG.Response(http.StatusOK, code, data)
}

//商品详情
func DetailGoods(c *gin.Context) {
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
	var data *models.Goods

	//获取商品详情
	goodsService := goods.Goods{
		Id: id,
	}

	goodsRst, err := goodsService.DetailGoods()
	if err.Code == 0 {
		code = e.SUCCESS
		data = goodsRst
	}

	appG.Response(http.StatusOK, code, data)

}

//删除文章
func DeleteGoods(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	//获取商品详情
	goodsService := goods.Goods{
		Id: id,
	}

	if !valid.HasErrors() {
		selectGoods, selectErr := goodsService.DetailGoods() //判断商品是否存在
		deleteErr := goodsService.DeleteGoods()              //删除操作

		if selectErr.Code == 0 && !deleteErr && selectGoods != nil {
			code = e.SUCCESS
		} else {
			code = e.ERROR_SQL_FAIL
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(fmt.Sprintf("%s,%s", "err key is "+err.Key, "err Message is "+err.Message))
		}
	}
	appG.Response(http.StatusOK, code, nil)

}

func EditGoods(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.Param("id")).MustInt()
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
	valid.Required(id, "id").Message("商品ID不能为空")
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
		err := goodsService.EditGoods(id)

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
