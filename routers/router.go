package routers

import (
	"gitee.com/muzipp/Distribution/middleware/auth"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/routers/admin"
	"gitee.com/muzipp/Distribution/routers/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(cors.Default())

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.Server{}.RunMode)

	//设置鉴权路由（就是登录接口）
	r.POST("/auth", common.GetAuth)

	/**
	创建一个路由组，路由组的路由可以具有相同的路由前缀或者中间件
	*/
	apiAdmin := r.Group("api/admin")

	//使用鉴权中间件(后面的所有的api/admin开头的路由都会经过这个中间件)
	apiAdmin.Use(auth.Auth())

	//{}代表作用域，作用内的变量只在作用域内有效，作用域内的a变量是没法在作用域外访问的
	{
		apiAdmin.POST("/logout", common.Logout)

		//会员
		apiAdmin.POST("/member", admin.AddMember)
		apiAdmin.GET("/member", admin.ListMembers)
		apiAdmin.POST("/member/statusChange", admin.MemberStatusChange)
		apiAdmin.PUT("/member/:id", admin.EditMember)
		//apiAdmin.GET("/member/:id", admin.DetailMember)

		//商品
		apiAdmin.POST("/goods", admin.AddGoods)
		apiAdmin.GET("/goods", admin.ListGoods)
		apiAdmin.GET("/goods/:id", admin.DetailGoods)
		apiAdmin.DELETE("/goods/:id", admin.DeleteGoods)

	}

	return r
}
