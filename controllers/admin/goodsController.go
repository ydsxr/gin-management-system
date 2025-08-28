package admin

import (
	"fmt"
	"go_demo/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoodsController struct {
	BaseController
}

func (con GoodsController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{})
}
func (con GoodsController) Add(c *gin.Context) {
	// 获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=?", 0).Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("goods_cate.sort ASC")
	}).Order("sort ASC").Find(&goodsCateList)
	// 获取商品颜色信息
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	// 获取商品规格包装
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goods/add.html", gin.H{
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
	})
}
func (con GoodsController) ImageUpload(c *gin.Context) {
	imageUrl, err := models.UploadImg(c, "file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"link": "/" + imageUrl,
	})
}
func (con GoodsController) GoodsTypeAttribute(c *gin.Context) {
	cateId, err1 := models.Int(c.Query("cateId"))
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err2 := models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttributeList).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"result":  "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  goodsTypeAttributeList,
		})
	}
}
func (con GoodsController) DoAdd(c *gin.Context) {
	// 表单中有多个相同的name
	attrIdList := c.PostFormArray("attr_id_list")
	attrValueList := c.PostFormArray("attr_value_list")
	goodsImageList := c.PostFormArray("goods_image_list")
	c.JSON(200,gin.H{
		"attrIdList":attrIdList,
		"attrValueList":attrValueList,
		"goodsImageList":goodsImageList,
	})

}
func (con GoodsController) Edit(c *gin.Context) {
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

	c.HTML(http.StatusOK, "admin/goods/edit.html", gin.H{
		"goodsCateList": goodsCateList,
		"goodsCate":     goodsCate,
	})

}
func (con GoodsController) DoEdit(c *gin.Context) {
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
		con.Error(c, "修改失败", "/admin/goods/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改成功", "/admin/goodsCate")

}
func (con GoodsController) Delete(c *gin.Context) {
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
