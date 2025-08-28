package api

import "github.com/gin-gonic/gin"

type ApiController struct{}

func (con ApiController) Index(c *gin.Context) {
	c.String(200, "我是一个api接口")
}
func (con ApiController) Userlist(c *gin.Context) {
	c.String(200, "我是一个api接口-Userlist")
}
func (con ApiController) Plist(c *gin.Context) {
	c.String(200, "我是一个api接口-Plist")
}
