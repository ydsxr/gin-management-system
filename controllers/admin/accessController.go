package admin

import (
	"fmt"
	"go_demo/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessController struct {
	BaseController
}

func (con AccessController) Index(c *gin.Context) {
	accessList := []models.Access{}
	// 自关联
	models.DB.Where("module_id=?", 0).Preload("AccessItem").Find(&accessList)
	fmt.Printf("%#v", accessList)
	c.HTML(200, "admin/access/index.html", gin.H{
		"accessList": accessList,
	})
}
func (con AccessController) Add(c *gin.Context) {
	// 获取顶级模块
	accessList := []models.Access{}
	// 获取顶级分类
	models.DB.Where("module_id=?", 0).Find(&accessList)
	c.HTML(200, "admin/access/add.html", gin.H{
		"accessList": accessList,
	})
}
func (con AccessController) Edit(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "参数错误", "/admin/access")
		return
	}
	access := models.Access{Id: id}
	models.DB.Find(&access)

	// 获取顶级模块
	accessList := []models.Access{}
	// 获取顶级分类
	models.DB.Where("module_id=?", 0).Find(&accessList)
	c.HTML(200, "admin/access/edit.html", gin.H{
		"accessList": accessList,
		"access":     access,
	})
}
func (con AccessController) DoEdit(c *gin.Context) {
	// 获取表单信息
	// id用于错误回溯
	id, err := models.Int(c.PostForm("id"))
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	accessType, err1 := models.Int(c.PostForm("type"))
	actionName := c.PostForm("action_name")
	url := c.PostForm("url")
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	description := c.PostForm("description")
	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/access")
		return
	}
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", "/admin/access/edit?id="+models.String(id))
		return
	}
	access := models.Access{
		Id: id,
	}
	models.DB.Find(&access)
	access.ModuleName = moduleName
	access.Type = accessType
	access.ActionName = actionName
	access.Url = url
	access.ModuleId = moduleId
	access.Sort = sort
	access.Description = description
	access.Status = status
	// 保存修改后的数据
	err5 := models.DB.Save(&access).Error
	if err5 != nil {
		con.Error(c, "修改权限失败", "/admin/access/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改权限成功", "/admin/access")
}
func (con AccessController) DoAdd(c *gin.Context) {
	// 获取表单信息
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	accessType, err1 := models.Int(c.PostForm("type"))
	actionName := c.PostForm("action_name")
	url := c.PostForm("url")
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	description := c.PostForm("description")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/access/add")
		return
	}
	access := models.Access{
		ModuleName:  moduleName,
		Type:        accessType,
		ActionName:  actionName,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}
	err5 := models.DB.Create(&access).Error
	if err5 != nil {
		con.Error(c, "增加权限失败", "/admin/access/add")
		return
	}
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", "/admin/access/add")
		return
	}
	con.Success(c, "增加权限成功", "/admin/access")
}
func (con AccessController) Delete(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/access")
	} else {
		access := models.Access{Id: id}
		models.DB.Find(&access)
		// 先判断是否是顶级模块
		if access.ModuleId == 0 { //是顶级模块
			// 判断是否有子模块
			accessList := []models.Access{}
			models.DB.Where("module_id=?", access.Id).Find(&accessList)
			if len(accessList) > 0 {
				con.Error(c, "当前模块下有子模块，删除错误", "/admin/access")
			} else { //没有子模块
				err2 := models.DB.Delete(&access).Error
				if err2 != nil {
					con.Error(c, "删除失败", "/admin/access")
					return
				}
				con.Success(c, "删除成功", "/admin/access")
			}
		} else { //不是顶级模块
			err3 := models.DB.Delete(&access).Error
			if err3 != nil {
				con.Error(c, "删除失败", "/admin/access")
				return
			}
			con.Success(c, "删除成功", "/admin/access")
		}

	}
}
