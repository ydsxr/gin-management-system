package admin

import (
	"encoding/json"
	"fmt"
	"go_demo/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	// fmt.Println(models.MD5("123456"))
	c.HTML(200, "admin/login/login.html", gin.H{})
}
func (con LoginController) DoLogin(c *gin.Context) {
	// 和html中的name对应
	username := c.PostForm("username")
	password := c.PostForm("password")

	captchaId := c.PostForm("captchaId")
	verifyValue := c.PostForm("verifyValue")

	// 1，验证验证码是否正确
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		// 2，查询数据库，验证用户密码登录
		userinfo := []models.Manager{}
		password = models.MD5(password)
		fmt.Println("username:", username)
		fmt.Println("password:", password)
	    models.DB.Where("username = ? AND password = ?", username, password).Find(&userinfo)
		// c.JSON(http.StatusOK,userinfo)
		if len(userinfo) > 0 {
			session := sessions.Default(c)
			// set只能保存字符串，不能保存结构体切片，需要先转换成Json字符串
			userinfoSlice, _ := json.Marshal(userinfo)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")
		} else {
			con.Error(c, "用户名或密码错误", "/admin/login")
		}
	} else {
		con.Error(c, "验证码验证失败", "/admin/login")
	}
}
func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, _, err := models.MakeCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
	// c.HTML(200, "admin/login/login.html", gin.H{})
}

func (con LoginController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	con.Success(c, "退出登录成功", "/admin/login")
}
