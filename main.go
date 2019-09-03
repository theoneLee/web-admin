package main

import (
	"fmt"
	"github.com/gin-blog/models"
	"github.com/gin-blog/pkg/gredis"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/pkg/setting"
	"github.com/gin-blog/routers"
	"log"
	"net/http"
)

func main() {
	//注册所有路由
	setting.Setup()
	models.Setup()
	logging.Setup()
	err := gredis.Setup() //注册redis配置
	if err != nil {
		log.Printf("Service Start Faild: %s\n", err)
	}
	router := routers.InitRouter()

	//实例化一个服务器（地址/端口号/读取超时/写入超时/header头最大字节数）
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
