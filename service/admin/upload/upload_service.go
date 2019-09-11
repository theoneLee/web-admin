package upload

import (
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/upload"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

//上传商品图片
func Image(image []*multipart.FileHeader) (images []string, err e.SelfError) {

	//添加商品基础信息成功的情况，上传图片，记录对应的商品图片信息
	imageUrl, imageErr := uploadImages(image)

	if imageErr.Code != 0 { //图片上传失败
		err.Code = imageErr.Code
		return
	}

	//遍历图片
	for _, value := range imageUrl {
		images = append(images, value)
	}
	return
}

func uploadImages(images []*multipart.FileHeader) (imageUrl []string, err e.SelfError) {
	//获取上传图片的图片名称
	var c *gin.Context

	for _, image := range images {
		imageName := upload.GetImageName(image.Filename)

		//获取图片完整路径
		fullPath := upload.GetImageFullPath()

		//获取图片保存路径
		savePath := upload.GetImagePath()

		//获取完整路径+文件名
		src := fullPath + imageName

		//检测文件后缀和文件大小
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(image.Size) {
			err.Code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			//检测图片
			imageError := upload.CheckImage(fullPath)
			if imageError != nil {
				err.Code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if imageSaveError := c.SaveUploadedFile(image, src); imageSaveError != nil {
				err.Code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				imageUrl = append(imageUrl, savePath+imageName)
			}
		}
	}

	return
}
