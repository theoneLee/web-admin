package admin

import (
	"fmt"
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
	levelId := com.StrTo(c.DefaultPostForm("level_id", "1")).MustInt()
	name := c.DefaultPostForm("name", "")
	idCard := c.DefaultPostForm("id_card", "")
	birth := c.DefaultPostForm("birth", "")
	phone := c.DefaultPostForm("phone", "")
	sparePhone := c.DefaultPostForm("spare_phone", "")
	email := c.DefaultPostForm("email", "")
	bankCard := c.DefaultPostForm("bank_card", "")
	bank := c.DefaultPostForm("bank", "")
	relationId := com.StrTo(c.DefaultPostForm("relation_id", "0")).MustInt()

	valid := validation.Validation{}
	valid.Required(sex, "sex").Message("性别不能为空")
	valid.Required(name, "name").Message("姓名不能为空")
	valid.Required(idCard, "id_card").Message("身份证不能为空")
	valid.Required(birth, "birth").Message("生日不能为空")
	valid.Required(phone, "phone").Message("电话地址不能为空")
	valid.Required(sparePhone, "spare_phone").Message("备用电话不能为空")
	valid.Required(email, "email").Message("邮箱不能为空")
	valid.Required(bankCard, "bank_card").Message("银行卡号不能为空")
	valid.Required(bank, "bank").Message("开户行不能为空")

	//设置返回数据
	data := make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		memberService := member.Member{
			Sex:        sex,
			LevelId:    levelId,
			Name:       name,
			IdCard:     idCard,
			Birth:      birth,
			Phone:      phone,
			SparePhone: sparePhone,
			Email:      email,
			BankCard:   bankCard,
			Bank:       bank,
			Status:     1,
			RelationId: relationId,
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
			Id:id,
			Status:status,
		}
		code = e.ERROR_SQL_FAIL
		err := memberService.StatusChange()
		if err.Code == 0  {
			code = e.SUCCESS
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}


	appG.Response(http.StatusOK, code, nil)

}

//
////会员详情
//func DetailMember(c *gin.Context) {
//	appG := app.Gin{C: c} //实例化响应对象
//
//}
