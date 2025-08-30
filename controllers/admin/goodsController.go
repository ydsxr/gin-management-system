package admin

import (
	"fmt"
	"go_demo/models"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 定义一个协程标记
var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

func (con GoodsController) Index(c *gin.Context) {
	goodsList := []models.Goods{}
	models.DB.Find(&goodsList)
	c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
		"goodsList": goodsList,
	})
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
	// 记录goods的表有三张：goods、goods_attr、goods_image，删除数据要同步删除
	//1、获取表单提交过来的数据
	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	// goodsSn := c.PostForm("goods_sn")
	cateId, err1 := models.Int(c.PostForm("cate_id"))
	//注意小数点
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))
	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color")
	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	//isDelete, err5 := models.Int(c.PostForm("is_delete"))//
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, err2 := models.Int(c.PostForm("goods_type_id"))
	sort, err10 := models.Int(c.PostForm("sort"))
	goodsNumber, err := models.Int(c.PostForm("goods_number"))
	status, err3 := models.Int(c.PostForm("status"))
	addTime := int(models.GetUnix())
	if err1 != nil || err2 != nil || err3 != nil || err10 != nil || err != nil {
		con.Error(c, "获取数据失败", "/admin/goods/add")
		return
	}
	//2、获取颜色信息 把颜色转化为字符串
	goodsColorStr := strings.Join(goodsColorArr, ",") //把切片里的放在一起

	//3、上传图片 生成缩略图
	goodsImg, _ := models.UploadImg(c, "goods_img")
	//4、增加商品数据
	goods := models.Goods{
		Title:    title,
		SubTitle: subTitle,
		// GoodsSn:       goodsSn,
		CateId:        cateId,
		ClickCount:    100,
		GoodsNumber:   goodsNumber,
		MarketPrice:   marketPrice,
		Price:         price,
		RelationGoods: relationGoods,
		GoodsAttr:     goodsAttr,
		GoodsVersion:  goodsVersion,
		GoodsGift:     goodsGift,
		GoodsFitting:  goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc:     goodsDesc,
		GoodsContent:  goodsContent,
		//IsDelete:      isDelete,
		IsHot:       isHot,
		IsBest:      isBest,
		IsNew:       isNew,
		GoodsTypeId: goodsTypeId,
		Sort:        sort,
		Status:      status,
		AddTime:     addTime,
		GoodsColor:  goodsColorStr,
		GoodsImg:    goodsImg,
	}
	err4 := models.DB.Create(&goods).Error
	if err4 != nil {
		con.Error(c, "增加商品失败", "/admin/goods/add")
		return
	}
	//5、增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done() // 标记减一
	}()
	//6、增加规格包装 attrIdList和attrValueList一一对应
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, errAttributeId := models.Int(attrIdList[i])
			fmt.Println("商品属性Id:", goodsTypeAttributeId)
			if errAttributeId != nil {
				con.Error(c, "获取商品类型错误", "/admin/goods/add")
				return
			}
			// 获取商品类型属性的数据
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)

			// 给商品属性里面增加数据 规格包装
			goodsAttrObj := models.GoodsAttr{}
			goodsAttrObj.GoodsId = goods.Id
			goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
			goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
			goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
			goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
			goodsAttrObj.AttributeValue = attrValueList[i]
			goodsAttrObj.Status = 1
			goodsAttrObj.Sort = 10
			goodsAttrObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsAttrObj)
		}
		wg.Done() // 标记减一
	}()
	wg.Wait() // 标记为0时开始往下执行
	con.Success(c, "增加商品成功", "/admin/goods")

}
func (con GoodsController) Edit(c *gin.Context) {
	//1、获取要修改商品数据
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goods")
	}
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	//2、获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=?", 0).Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("goods_cate.sort ASC")
	}).Order("sort ASC").Find(&goodsCateList)
	//3、获取所有颜色 以及选中的颜色
	goodsColorSlice := strings.Split(goods.GoodsColor, ",") //切片 类似于数组 元素可以不是数
	goodsColorMap := make(map[string]string)                // 转化为map类型是为了做对比
	for _, v := range goodsColorSlice {
		goodsColorMap[v] = v
	}
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	for i := 0; i < len(goodsColorList); i++ {
		_, ok := goodsColorMap[models.String(goodsColorList[i].Id)]
		if ok {
			goodsColorList[i].Checked = true
		}
	}
	//4、商品的图库信息
	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsImageList)
	//5、获取商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	//6、获取规格信息
	goodsAttrList := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsAttrList)
	goodsAttrStr := ""
	for _, v := range goodsAttrList {
		switch v.AttributeType {
		case 1:
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span> <input type="hidden" name="attr_id_list" value="%v" />   <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		case 2:
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span><input type="hidden" name="attr_id_list" value="%v" />  <textarea cols="50" rows="3" name="attr_value_list">%v</textarea></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		default:
			//获取当前类型对应的值
			goodsTypeArttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeArttribute)
			attrValueSlice := strings.Split(goodsTypeArttribute.AttrValue, "\n")

			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span>  <input type="hidden" name="attr_id_list" value="%v" /> `, v.AttributeTitle, v.AttributeId)
			goodsAttrStr += `<select name="attr_value_list">`
			for i := 0; i < len(attrValueSlice); i++ {
				if attrValueSlice[i] == v.AttributeValue {
					goodsAttrStr += fmt.Sprintf(`<option value="%v" selected >%v</option>`, attrValueSlice[i], attrValueSlice[i])
				} else {
					goodsAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice[i], attrValueSlice[i])
				}
			}
			goodsAttrStr += `</select>`
			goodsAttrStr += `</li>`
		}
	}

	// 获取商品规格包装

	c.HTML(http.StatusOK, "admin/goods/edit.html", gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
		"goodsAttrStr":   goodsAttrStr,
		"goodsImageList": goodsImageList,
	})
}
func (con GoodsController) DoEdit(c *gin.Context) {
	// 记录goods的表有三张：goods、goods_attr、goods_image，删除数据要同步删除
	//1、获取表单提交过来的数据
	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "获取ID失败", "/admin/goods")
		return
	}
	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	// goodsSn := c.PostForm("goods_sn")
	cateId, err2 := models.Int(c.PostForm("cate_id"))
	//注意小数点
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))
	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color")
	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	//isDelete, err5 := models.Int(c.PostForm("is_delete"))//
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, err3 := models.Int(c.PostForm("goods_type_id"))
	sort, err4 := models.Int(c.PostForm("sort"))
	goodsNumber, err5 := models.Int(c.PostForm("goods_number"))
	status, err6 := models.Int(c.PostForm("status"))
	addTime := int(models.GetUnix())
	if err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		con.Error(c, "获取数据失败", "/admin/goods/add")
		return
	}
	//2、获取颜色信息 把颜色转化为字符串
	goodsColorStr := strings.Join(goodsColorArr, ",") //把切片里的放在一起

	//3、修改数据
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	goods.Title = title
	goods.SubTitle = subTitle
	// GoodsSn:       goodsSn,
	goods.CateId = cateId
	goods.ClickCount = 100
	goods.GoodsNumber = goodsNumber
	goods.MarketPrice = marketPrice
	goods.Price = price
	goods.RelationGoods = relationGoods
	goods.GoodsAttr = goodsAttr
	goods.GoodsVersion = goodsVersion
	goods.GoodsGift = goodsGift
	goods.GoodsFitting = goodsFitting
	goods.GoodsKeywords = goodsKeywords
	goods.GoodsDesc = goodsDesc
	goods.GoodsContent = goodsContent
	//IsDelete:      isDelete,
	goods.IsHot = isHot
	goods.IsBest = isBest
	goods.IsNew = isNew
	goods.GoodsTypeId = goodsTypeId
	goods.Sort = sort
	goods.Status = status
	goods.AddTime = addTime
	goods.GoodsColor = goodsColorStr

	//4、上传图片 生成缩略图
	goodsImg, err7 := models.UploadImg(c, "goods_img")
	if err7 == nil && len(goodsImg) > 0 {
		goods.GoodsImg = goodsImg
	}
	//5、执行修改
	err8 := models.DB.Save(&goods).Error
	if err8 != nil {
		con.Error(c, "修改商品失败", "/admin/goods/edit?id="+models.String(id))
		return
	}
	//6、修改图库信息 增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done() // 标记减一
	}()
	//7、修改规格包装 attrIdList和attrValueList一一对应
	//删除当前商品下面的规格包装
	goodsAttrObj := models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Delete(&goodsAttrObj)
	//重新执行增加
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, errAttributeId := models.Int(attrIdList[i])
			fmt.Println("商品属性Id:", goodsTypeAttributeId)
			if errAttributeId != nil {
				con.Error(c, "获取商品类型错误", "/admin/goods/add")
				return
			}
			// 获取商品类型属性的数据
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)

			// 给商品属性里面增加数据 规格包装
			goodsAttrObj := models.GoodsAttr{}
			goodsAttrObj.GoodsId = goods.Id
			goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
			goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
			goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
			goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
			goodsAttrObj.AttributeValue = attrValueList[i]
			goodsAttrObj.Status = 1
			goodsAttrObj.Sort = 10
			goodsAttrObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsAttrObj)
		}
		wg.Done() // 标记减一
	}()
	wg.Wait() // 标记为0时开始往下执行
	con.Success(c, "修改商品成功", "/admin/goods")

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
