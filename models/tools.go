package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 时间戳转换成日期
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 日期转换成时间戳 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前的日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// MD5加密
func MD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//把字符串解析成html
func Str2Html(str string) template.HTML{
	return template.HTML(str)
}

// 转string为int
func Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// 转string为Float64
func Float(str string) (float64, error) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err
}

// 转int为string
func String(n int) string {
	str := strconv.Itoa(n)
	return str
}

// 上传图片
func UploadImg(c *gin.Context, picName string) (string, error) {
	//1、获取上传的文件
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}
	//2、获取文件后缀名
	extName := path.Ext(file.Filename)
	allowExtName := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtName[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}
	//3、创建文件保存目录
	day := GetDay()
	dir := "./static/updata/" + day //路径
	err1 := os.MkdirAll(dir, 0666)  //0666代表目录权限
	if err1 != nil {
		return "", err1
	}
	//4、生成文件名称和保存的目录
	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName
	//5、执行上传
	dst := path.Join(dir, fileName) //打包文件名和路径
	c.SaveUploadedFile(file, dst)
	return dst, nil
}
