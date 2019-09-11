package admin

import (
	"gitee.com/muzipp/Distribution/pkg/app"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/service/admin/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

//添加商品
func UploadGoodsImage(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	c.PostForm("image")
	image := c.Request.MultipartForm.File["image"]

	//设置返回数据
	data := make([]string, 5)

	code := e.INVALID_PARAMS
	//图片上传
	images, err := upload.Image(image)

	if err.Code == 0 {
		code = e.SUCCESS
		data = images
	}

	appG.Response(http.StatusOK, code, data)
}
