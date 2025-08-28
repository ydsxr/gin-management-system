package admin

import (
	"go_demo/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoodsTypeAttributeController struct {
	BaseController
}

func (con GoodsTypeAttributeController) Index(c *gin.Context) {
	cateId, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
	}
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	models.DB.Where("cate_id", cateId).Find(&goodsTypeAttributeList)

	// 获取对应的类型
	goodsType := models.GoodsType{}
	models.DB.Where("id=?", cateId).Find(&goodsType)

	c.HTML(200, "admin/goodsTypeAttribute/index.html", gin.H{
		"cateId":                 cateId,
		"goodsTypeAttributeList": goodsTypeAttributeList,
		"goodsType":              goodsType,
	})
}
func (con GoodsTypeAttributeController) Add(c *gin.Context) {
	// 获取当前商品类型属性对应的商品类型id
	cateId, err1 := models.Int(c.Query("cate_id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
		return
	}
	// 获取所有的商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	c.HTML(200, "admin/goodsTypeAttribute/add.html", gin.H{
		"goodsTypeList": goodsTypeList,
		"cateId":        cateId,
	})
}
func (con GoodsTypeAttributeController) DoAdd(c *gin.Context) {
	// trim去除表单传过来的title中的空格
	title := strings.Trim(c.PostForm("title"), " ")
	cateId, err1 := models.Int(c.PostForm("cate_id"))
	attrType, err2 := models.Int(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err3 := models.Int(c.PostForm("sort"))
	if err1 != nil || err2 != nil {
		con.Error(c, "非法请求", "/admin/goodsType")
		return
	}
	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	if err3 != nil {
		con.Error(c, "排序值不对", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}

	goodsTypeAttr := models.GoodsTypeAttribute{
		Title:     title,
		CateId:    cateId,
		AttrType:  attrType,
		AttrValue: attrValue,
		Status:    1,
		Sort:      sort,
		AddTime:   int(models.GetUnix()),
	}
	err := models.DB.Create(&goodsTypeAttr).Error
	if err != nil {
		con.Error(c, "增加商品类型属性失败 请重试", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
	} else {
		con.Success(c, "增加商品类型属性成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
	}
}

func (con GoodsTypeAttributeController) Edit(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
		return
	}
	// 获取所有类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	// 获取指定的商品属性
	goodsTyprAttribute := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTyprAttribute)

	c.HTML(200, "admin/goodsTypeAttribute/edit.html", gin.H{
		"goodsTyprAttribute": goodsTyprAttribute,
		"goodsTypeList":      goodsTypeList,
	})
}
func (con GoodsTypeAttributeController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	// trim去除表单传过来的title中的空格
	title := strings.Trim(c.PostForm("title"), " ")
	cateId, err2 := models.Int(c.PostForm("cate_id"))
	attrType, err3 := models.Int(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err4 := models.Int(c.PostForm("sort"))
	if err1 != nil || err2 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/goodsType")
		return
	}
	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	if err4 != nil {
		con.Error(c, "排序值不对", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}

	goodsTypeAttr := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttr)
	goodsTypeAttr.Title = title
	goodsTypeAttr.CateId = cateId
	goodsTypeAttr.AttrType = attrType
	goodsTypeAttr.AttrValue = attrValue
	goodsTypeAttr.Sort = sort

	err := models.DB.Save(&goodsTypeAttr).Error
	if err != nil {
		con.Error(c, "修改商品类型属性失败 请重试", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
	} else {
		con.Success(c, "修改商品类型属性成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
	}
}
func (con GoodsTypeAttributeController) Delete(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	cateId,err2:=models.Int(c.Query("cate_id"))
	if err1 != nil||err2!=nil {
		con.Error(c, "传入数据错误", "/admin/goodsType")
	} else {
		goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
		models.DB.Delete(&goodsTypeAttribute)
		con.Success(c, "删除该商品属性成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
	}
}
