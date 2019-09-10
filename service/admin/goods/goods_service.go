package goods

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	seleGoods "gitee.com/muzipp/Distribution/pkg/goods"
	"gitee.com/muzipp/Distribution/pkg/upload"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strings"
)

type Goods struct {
	Id     int
	Name   string
	Stock  int
	Price  float64
	Remark string
	Status int
	Images []string
	Offset int
	Limit  int
}

func (g *Goods) AddGoods(image []*multipart.FileHeader) (err e.SelfError) {
	data := make(map[string]interface{})
	data["name"] = g.Name
	data["remark"] = g.Remark
	data["price"] = g.Price
	data["stock"] = g.Stock
	data["status"] = g.Status

	//开启事务
	tx := models.Db.Begin()
	id, res := models.AddGoods(data, tx) //添加商品

	if res { //添加商品失败
		tx.Rollback() //事务回滚
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	//添加商品基础信息成功的情况，上传图片，记录对应的商品图片信息
	imageUrl, imageErr := uploadImages(image)

	if imageErr.Code != 0 { //图片上传失败
		tx.Rollback() //事务回滚
		err.Code = imageErr.Code
		return
	}

	//遍历图片
	dataImage := make(map[string]interface{})
	for _, value := range imageUrl {
		dataImage["goods_id"] = id
		dataImage["img"] = value
		res := models.AddGoodsImg(dataImage, tx)
		if res { //添加商品图片失败
			tx.Rollback() //事务回滚
			err.Code = e.ERROR_SQL_FAIL
			return
		}
	}

	tx.Commit() //事务提交
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

//商品列表
func (g *Goods) ListGoods() (goods []models.Goods, err e.SelfError) {
	fields := "*"
	goods, goodsErr := models.ListGoods(g.Offset, g.Limit, g.getMaps(), fields)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	//切割字符串图片
	for key, value := range goods {
		goods[key].Images = strings.Fields(value.Img)
		goods[key].StatusDesc = seleGoods.GetStatus(value.Status)
	}

	return

}

//商品数量
func (g *Goods) CountGoods() (count int, err e.SelfError) {
	count, goodsErr := models.CountGoods(g.getMaps())
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

//获取文章（redis不存在读取数据库）
func (g *Goods) DetailGoods() (goods *models.Goods, err e.SelfError) {
	fields := "*"
	goods, goodsErr := models.DetailGoods(g.Id, fields)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	if goods != nil {
		goods.Images = strings.Fields(goods.Img)
		goods.StatusDesc = seleGoods.GetStatus(goods.Status)
	}
	return
}

//判断商品是否存在
func (g *Goods) DeleteGoods() (flag bool) {

	tx := models.Db.Begin()

	goodsErr := models.DeleteGoods(g.Id, tx)
	goodsImgErr := models.DeleteGoodsImg(g.Id, tx)

	if goodsErr || goodsImgErr {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return goodsErr || goodsImgErr
}

//封装搜索条件
func (g *Goods) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_at"] = 0
	return maps
}
