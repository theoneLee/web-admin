package routers

import (
	"gitee.com/muzipp/Distribution/pkg/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.Server{}.RunMode)

	return r
}
