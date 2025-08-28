package admin

import (
	"go_demo/models"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	BaseController
}

func (con PersonController) Index(c *gin.Context) {
	c.String(200, "我是基本信息")
}
func (con PersonController) Resume(c *gin.Context) {
	person := models.Person{}
	models.DB.Where("name=?", "yds").Find(&person)
	c.HTML(200, "admin/person/resume.html", gin.H{
		"person": person,
	})
}
