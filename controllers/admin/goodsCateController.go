package admin

import (
	"fmt"
	"go_demo/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type GoodsCateController struct {
	BaseController
}

func (con GoodsCateController) Index(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Order("sort ASC").Find(&goodsCateList)
	c.HTML(http.StatusOK, "admin/goodsCate/index.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}
func (con GoodsCateController) Add(c *gin.Context) {
	// 获取顶级分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Order("sort ASC").Find(&goodsCateList)

	c.HTML(http.StatusOK, "admin/goodsCate/add.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}
func (con GoodsCateController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	pid, err1 := models.Int(c.PostForm("pid"))
	link := c.PostForm("link")
	template := c.PostForm("template")
	subTitle := c.PostForm("subTitle")
	keywords := c.PostForm("keywords")
	sort, err2 := models.Int(c.PostForm("sort"))
	status, err3 := models.Int(c.PostForm("status"))
	if err1 != nil || err3 != nil {
		con.Error(c, "传入参数错误", "admin/goodsCate/add")
		return
	}
	if err2 != nil {
		con.Error(c, "排序值必须为整数", "admin/goodsCate/add")
		return
	}
	cateImg, err4 := models.UploadImg(c, "cate_img")
	if err4 != nil {
		fmt.Println(err4)
	}
	goodsCate := models.GoodsCate{
		Title:    title,
		Pid:      pid,
		SubTitle: subTitle,
		Link:     link,
		Template: template,
		Keywords: keywords,
		Sort:     sort,
		Status:   status,
		CateImg:  cateImg,
		AddTime:  int(models.GetUnix()),
	}
	err5 := models.DB.Create(&goodsCate).Error
	if err5 != nil {
		con.Error(c, "增加数据失败", "/admin/goodsCate/add")
		return
	}
	con.Success(c, "增加数据成功", "/admin/goodsCate")

}
func (con GoodsCateController) Edit(c *gin.Context) {
	// 获取商品分类Id
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "Id获取错误", "/admin/goodsCate")
		return
	}

	// 获取顶级分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Order("sort ASC").Find(&goodsCateList)

	// 获取该商品数据
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)

	c.HTML(http.StatusOK, "admin/goodsCate/edit.html", gin.H{
		"goodsCateList": goodsCateList,
		"goodsCate":     goodsCate,
	})

}
func (con GoodsCateController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	title := c.PostForm("title")
	link := c.PostForm("link")
	template := c.PostForm("template")
	pid, err2 := models.Int(c.PostForm("pid"))
	subTitle := c.PostForm("sub_title")
	keywords := c.PostForm("keywords")
	description := c.PostForm("description")
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	if err1 != nil || err2 != nil || err4 != nil {
		con.Error(c, "参数类型不正确", "/admin/goodsCate")
		return
	}
	if err3 != nil {
		con.Error(c, "排序值需要为整数", "/admin/goodsCate")
	}
	cateImg, err5 := models.UploadImg(c, "cate_img")
	if err5 != nil {
		fmt.Println(err5)
	}
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	goodsCate.Title = title
	goodsCate.Link = link
	goodsCate.Template = template
	goodsCate.Pid = pid
	goodsCate.SubTitle = subTitle
	goodsCate.Keywords = keywords
	goodsCate.Description = description
	goodsCate.Status = status
	goodsCate.Sort = sort
	if cateImg != "" {
		goodsCate.CateImg = cateImg
	}
	err6 := models.DB.Save(&goodsCate).Error
	if err6 != nil {
		con.Error(c, "修改失败", "/admin/goodsCate/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改成功", "/admin/goodsCate")

}
func (con GoodsCateController) Delete(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsCate")
	}
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	ImgSrc := goodsCate.CateImg
	err2 := models.DB.Delete(&goodsCate).Error
	if err2 != nil {
		con.Error(c, "删除数据失败", "/admin/goodsCate")
		return
	}
	err3 := os.Remove("./" + ImgSrc)
	if err3 != nil {
		con.Error(c, "删除图片失败", "/admin/goodsCate")
	}
	con.Success(c, "删除数据成功", "/admin/goodsCate")
}
