package upload

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"gitee.com/muzipp/Distribution/pkg/file"
	"gitee.com/muzipp/Distribution/pkg/setting"
	"gitee.com/muzipp/Distribution/pkg/util"
)

//获取图片的完整URL
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

//获取图片名称+后缀
func GetImageName(name string) string {

	//获取文件后缀
	ext := path.Ext(name)

	//处理文件名（如果文件名不包括后缀的话，直接返回文件名），如果文件名包含后缀的话，处理文件名（取出指定的后缀）
	fileName := strings.TrimSuffix(name, ext)

	//文件名拼接上时间戳和随机数
	timeNow := time.Now().Format("20060102150405")
	tempRand := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z",
	}
	for i := 0; i <= 3; i++ {
		timeNow = timeNow + tempRand[rand.Intn(26)]
	}

	//文件名md5处理
	fileName = util.EncodeMD5(fileName+timeNow)

	//返回文件名+后缀
	return fileName + ext
}

//获取图片保存路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

//获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

//检测文件后缀是否符合配置设置的文件后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

//检测文件大小
func CheckImageSize(size int64) bool {
	return int(size) <= setting.AppSetting.ImageMaxSize
}

//检查图片
func CheckImage(src string) error {

	//获取当前目录所在的根目录
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	//判断文件目录是否存在，不存在创建目录
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	//检测权限
	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
