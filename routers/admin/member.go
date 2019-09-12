package admin

import (
	"fmt"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
	"gitee.com/muzipp/Distribution/service/admin/member"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

//添加会员
func AddMember(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	sex := com.StrTo(c.DefaultPostForm("sex", "0")).MustInt()
	levelId := com.StrTo(c.DefaultPostForm("level_id", "0")).MustInt()
	name := c.DefaultPostForm("name", "")
	relationName := c.DefaultPostForm("relation_name", "")
	idCard := c.DefaultPostForm("id_card", "")
	birth := c.DefaultPostForm("birth", "")
	phone := c.DefaultPostForm("phone", "")
	sparePhone := c.DefaultPostForm("spare_phone", "")
	email := c.DefaultPostForm("email", "")
	bankCard := c.DefaultPostForm("bank_card", "")
	bank := c.DefaultPostForm("bank", "")
	password := c.DefaultPostForm("password", "")
	username := c.DefaultPostForm("username", "")
	isOperate := com.StrTo(c.DefaultPostForm("is_operate", "0")).MustInt()
	operateAddress := c.DefaultPostForm("operate_address", "")
	remark := c.DefaultPostForm("remark", "")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("姓名不能为空")
	valid.Required(password, "password").Message("密码不能为空")
	valid.Required(username, "username").Message("用户名不能为空")
	valid.Required(levelId, "level_id").Message("等级不能为空")
	valid.Required(relationName, "relation_name").Message("上级代理不能为空")

	//设置返回数据
	data := make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {

		//校验用户名是否重复
		user := models.CheckAuth(username, 0)
		if user.ID > 0 {
			code = e.ERROR_USERNAME
			goto End
		}

		relationUser := models.CheckAuth(relationName, 0)
		if relationUser.ID == 0 {
			code = e.ERROR_USERNAME
			goto End
		}

		memberService := member.Member{
			Sex:            sex,
			LevelId:        levelId,
			Name:           name,
			IdCard:         idCard,
			Birth:          birth,
			Phone:          phone,
			SparePhone:     sparePhone,
			Email:          email,
			BankCard:       bankCard,
			Bank:           bank,
			Status:         1,
			RelationId:     relationUser.ID,
			Username:       username,
			PassWord:       util.EncodeMD5(password),
			IsOperate:      isOperate,
			OperateAddress: operateAddress,
			Remark:         remark,
		}
		err := memberService.AddMember()

		if err.Code == 0 {
			code = e.SUCCESS
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info(fmt.Sprintf("%s,%s", "err key is "+err.Key, "err Message is "+err.Message))
		}
	}

End:
	appG.Response(http.StatusOK, code, data)

}

//会员列表
func ListMembers(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	data := make(map[string]interface{})

	memberService := member.Member{
		Offset: util.GetPage(c),
		Limit:  setting.AppSetting.PageSize,
	}

	code := e.ERROR_SQL_FAIL
	members, err := memberService.ListMembers()
	total, totalError := memberService.CountMembers()

	if err.Code == 0 && totalError.Code == 0 {
		code = e.SUCCESS
		data["lists"] = members
		data["total"] = total
	}

	appG.Response(http.StatusOK, code, data)

}

//商品详情
func DetailMember(c *gin.Context) {
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
	var data *models.Member

	//获取商品详情
	goodsService := member.Member{
		Id: id,
	}

	memberRst, err := goodsService.DetailMember()
	if err.Code == 0 {
		code = e.SUCCESS
		data = memberRst
	}

	appG.Response(http.StatusOK, code, data)

}

//会员列表
func MemberStatusChange(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.PostForm("id")).MustInt()
	status := com.StrTo(c.PostForm("status")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Range(status, -2, -1, "status").Message("状态只允许-2或-1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		memberService := member.Member{
			Id:     id,
			Status: status,
		}
		code = e.ERROR_SQL_FAIL
		err := memberService.StatusChange()
		if err.Code == 0 {
			code = e.SUCCESS
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	appG.Response(http.StatusOK, code, nil)

}

func EditMember(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.Param("id")).MustInt()
	sex := com.StrTo(c.DefaultPostForm("sex", "0")).MustInt()
	levelId := com.StrTo(c.DefaultPostForm("level_id", "1")).MustInt()
	name := c.DefaultPostForm("name", "")
	idCard := c.DefaultPostForm("id_card", "")
	birth := c.DefaultPostForm("birth", "")
	phone := c.DefaultPostForm("phone", "")
	sparePhone := c.DefaultPostForm("spare_phone", "")
	email := c.DefaultPostForm("email", "")
	bankCard := c.DefaultPostForm("bank_card", "")
	bank := c.DefaultPostForm("bank", "")
	password := c.DefaultPostForm("password", "")
	isOperate := com.StrTo(c.DefaultPostForm("is_operate", "0")).MustInt()
	operateAddress := c.DefaultPostForm("operate_address", "")
	remark := c.DefaultPostForm("remark", "")

	valid := validation.Validation{}
	valid.Required(id, "id").Message("会员ID不能为空")
	valid.Required(name, "name").Message("姓名不能为空")
	valid.Required(password, "password").Message("密码不能为空")
	valid.Required(levelId, "level_id").Message("等级不能为空")

	//设置返回数据
	data := make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {

		//校验用户是否存在
		user := models.CheckUser(id)
		if user.ID == 0 {
			code = e.ERROR_USER
			goto End
		}

		memberService := member.Member{
			Sex:            sex,
			LevelId:        levelId,
			Name:           name,
			IdCard:         idCard,
			Birth:          birth,
			Phone:          phone,
			SparePhone:     sparePhone,
			Email:          email,
			BankCard:       bankCard,
			Bank:           bank,
			Status:         1,
			RelationId:     user.RelationId,
			PassWord:       password,
			IsOperate:      isOperate,
			OperateAddress: operateAddress,
			Remark:         remark,
		}

		err := memberService.EditMember(id)

		if err.Code == 0 {
			code = e.SUCCESS
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info(fmt.Sprintf("%s,%s", "err key is "+err.Key, "err Message is "+err.Message))
		}
	}

End:
	appG.Response(http.StatusOK, code, data)

}

