package itying

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// 加 . 调用函数时可以不用引入包名 不建议使用
	"github.com/hunterhug/go_image"
	qrcode "github.com/skip2/go-qrcode"
)

type DefaultController struct{}

func (con DefaultController) Thumbnail1(c *gin.Context) {
	// 按宽度进行比例缩放 输入输出都是文件
	filename := "static/updata/20250820/1755692223.jpg"
	savepath := "static/updata/20250820/0.jpg"
	err := go_image.ScaleF2F(filename, savepath, 600)
	if err != nil {
		c.String(http.StatusOK, "生成图片失败")
		return
	}
	c.String(http.StatusOK, "生成图片成功")

}
func (con DefaultController) GoQrcode(c *gin.Context) {
	var png []byte
	png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusOK, "生成二维码失败")
		return
	}
	c.String(http.StatusOK, string(png))

}
