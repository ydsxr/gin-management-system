package admin

import (
	"go_demo/models"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type NavController struct {
	BaseController
}

func (con NavController) Index(c *gin.Context) {
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	pageSize := 8
	navList := []models.Nav{}
	models.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&navList)
	// 获取总数量
	var count int64
	models.DB.Table("nav").Count(&count)
	c.HTML(200, "admin/nav/index.html", gin.H{
		"navList": navList,
		"page":page,
		"totalPages": math.Ceil(float64(float64(count) / float64(pageSize))),
	})
}
func (con NavController) Add(c *gin.Context) {
	c.HTML(200, "admin/nav/add.html", gin.H{})
}
func (con NavController) DoAdd(c *gin.Context) {
	// trim去除表单传过来的title中的空格
	title := strings.Trim(c.PostForm("title"), " ")
	if title == "" {
		con.Error(c, "导航名称不能为空", "/admin/nav/add")
		return
	}
	position, err1 := models.Int(c.PostForm("position"))
	relation := c.PostForm("relation")
	link := c.PostForm("link")
	isOpennew, err2 := models.Int(c.PostForm("is_opennew"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/add")
		return
	}
	nav := models.Nav{}
	nav.Title = title
	nav.Link = link
	nav.Position = position
	nav.Relation = relation
	nav.IsOpennew = isOpennew
	nav.Sort = sort
	nav.Status = status
	nav.AddTime = int(models.GetUnix())
	err := models.DB.Create(&nav).Error
	if err != nil {
		// 导航增加失败
		con.Error(c, "导航增加失败", "/admin/nav/add")
	} else {
		con.Success(c, "导航增加成功", "/admin/nav")
	}
}
func (con NavController) Edit(c *gin.Context) {
	// Query查询?后的数据
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/nav")
	} else {
		nav := models.Nav{Id: id}
		models.DB.Find(&nav)
		c.HTML(http.StatusOK, "admin/nav/edit.html", gin.H{
			"nav": nav,
		})
	}

}
func (con NavController) DoEdit(c *gin.Context) {
	//查询要修改的数据
	id, err1 := models.Int(c.PostForm("id"))
	title := strings.Trim(c.PostForm("title"), " ")
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/nav")
	} else {
		nav := models.Nav{Id: id}
		models.DB.Find(&nav)
		nav.Title = title
		err2 := models.DB.Save(&nav).Error
		if err2 != nil {
			con.Error(c, "修改导航信息失败", "/admin/nav/edit?id="+models.String(id))
		} else {
			con.Success(c, "修改导航信息成功", "/admin/nav")
		}
	}

}
func (con NavController) Delete(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/nav")
	} else {
		nav := models.Nav{Id: id}
		models.DB.Delete(&nav)
		con.Success(c, "删除导航成功", "/admin/nav")
	}
}
