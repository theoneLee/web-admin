package common

import (
	"encoding/json"
	"fmt"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/gredis"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"gitee.com/muzipp/Distribution/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type UserSession struct {
	Username string
	Name     string
	Id       int
}

var SelfToken string
var SelfUser UserSession

//登录接口
func GetAuth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}

	/**
	这种方式就直接验证了username和password必传，且最大长度为50
	*/
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		user := models.CheckAuth(username, 1)
		if user.ID > 0 && user.Password == util.EncodeMD5(password) { //用户存在且账号密码正确的情况

			//获取Redis中已经保存的token
			redisUserGetToken, _ := gredis.Get("user_id" + strconv.Itoa(user.ID))
			_ = json.Unmarshal(redisUserGetToken, &SelfToken)
			if SelfToken != "" {
				redisUserGetUser, _ := gredis.Get("user_token" + SelfToken)
				_ = json.Unmarshal(redisUserGetUser, &SelfUser)
				data["token"] = SelfToken
				data["name"] = SelfUser.Name
				code = e.SUCCESS
				goto End
			}

			/**
			生成token，生成的时候，同时设置了token有效期，
			*/
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				code = e.SUCCESS

				//token记录在Redis中，用来验证token登录
				redisTokenError := gredis.Set("user_token"+token, UserSession{Username: username, Id: user.ID, Name: user.Name}, 24*60*60)
				redisUserError := gredis.Set("user_id"+strconv.Itoa(user.ID), token, 24*60*60)

				if redisTokenError != nil || redisUserError != nil {
					code = e.ERROR_REDIS
					logging.Info(fmt.Sprintf("%s,%s", "redis token error is ", redisTokenError))
					logging.Info(fmt.Sprintf("%s,%s", "redis user error is ", redisUserError))
				} else {
					data["token"] = token
					data["name"] = user.Name
				}
			}

		} else {
			code = e.ERROR_AUTH
		}
	} else {

		/**
		验证出现错误的情况，遍历打印错误信息
		*/
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message) //记录错误日志
		}
	}

End:
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
