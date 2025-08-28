package admin

import (
	"fmt"
	"go_demo/models"
	"os"

	"github.com/gin-gonic/gin"
)

type FocusController struct {
	BaseController
}

func (con FocusController) Index(c *gin.Context) {
	focusList := []models.Focus{}
	models.DB.Find(&focusList)
	c.HTML(200, "admin/focus/index.html", gin.H{
		"focusList": focusList,
	})
}
func (con FocusController) Add(c *gin.Context) {
	c.HTML(200, "admin/focus/add.html", gin.H{})
}
func (con FocusController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	focusType, err1 := models.Int(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err2 := models.Int(c.PostForm("sort"))
	status, err3 := models.Int(c.PostForm("status"))
	if err1 != nil || err3 != nil {
		con.Error(c, "输入参数错误", "/admin/focus/add")
	}
	if err2 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/add")
	}
	// 上传文件
	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	if err4 != nil {
		fmt.Println(err4)
	}
	focus := models.Focus{
		Title:      title,
		FocusType:  focusType,
		FocusImage: focusImgSrc,
		Link:       link,
		Sort:       sort,
		Status:     status,
		AddTime:    int(models.GetUnix()),
	}
	err5 := models.DB.Create(&focus).Error
	if err5 != nil {
		con.Error(c, "增加轮播图失败", "/admin/focus/add")
		return
	}
	con.Success(c, "增加轮播图成功", "/admin/focus")
}
func (con FocusController) Edit(c *gin.Context) {
	focusId, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/focus")
	}
	focus := models.Focus{}
	models.DB.Where("id=?", focusId).Find(&focus)
	c.HTML(200, "admin/focus/edit.html", gin.H{
		"focus": focus,
	})
}
func (con FocusController) DoEdit(c *gin.Context) {
	focusId, err1 := models.Int(c.PostForm("id"))
	title := c.PostForm("title")
	focusType, err2 := models.Int(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	focusImgSrc, err5 := models.UploadImg(c, "focus_img")
	if err1 != nil {
		con.Error(c, "数据ID错误", "/admin/focus")
		return
	}
	//err2 != nil || err3 != nil || err4 != nil || err5 != nil
	if err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "数据传入失败", "/admin/focus/edit?id="+models.String(focusId))
		return
	}
	if err5 != nil {
		fmt.Println(err5)
	}
	focus := models.Focus{}
	models.DB.Where("id=?", focusId).Find(&focus)
	focus.Id = focusId
	focus.Title = title
	focus.FocusType = focusType
	focus.Link = link
	focus.Sort = sort
	focus.Status = status
	if focusImgSrc != "" {
		focus.FocusImage = focusImgSrc
	}
	err6 := models.DB.Save(&focus).Error
	if err6 != nil {
		con.Error(c, "修改轮播图失败", "/admin/focus/edit?id="+models.String(focusId))
		return
	}
	con.Success(c, "修改轮播图成功", "/admin/focus")

}
func (con FocusController) Delete(c *gin.Context) {
	focusId, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "获取ID失败", "/admin/focus")
		return
	}
	focus := models.Focus{Id: focusId}
	models.DB.Find(&focus)
	var src string = focus.FocusImage
	err2 := models.DB.Delete(&focus).Error
	if err2 != nil {
		con.Error(c, "删除数据失败", "/admin/focus")
		return
	}
	// 把服务器中的图片删了
	err3 := os.Remove("./" + src)
	if err3 != nil {
		con.Error(c, "删除图片失败", "/admin/focus")
		return
	}
	con.Success(c, "删除数据成功", "/admin/focus")
}
