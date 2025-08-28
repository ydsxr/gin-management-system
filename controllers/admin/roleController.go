package admin

import (
	"fmt"
	"go_demo/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	BaseController
}

func (con RoleController) Index(c *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(200, "admin/role/index.html", gin.H{
		"roleList": roleList,
	})
}
func (con RoleController) Add(c *gin.Context) {
	c.HTML(200, "admin/role/add.html", gin.H{})
}
func (con RoleController) DoAdd(c *gin.Context) {
	// trim去除表单传过来的title中的空格
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if title == "" {
		con.Error(c, "角色名称不能为空", "/admin/role/add")
		return
	}
	role := models.Role{}
	role.Title = title
	role.Description = description
	role.Status = 1
	role.AddTime = int(models.GetUnix())
	err := models.DB.Create(&role).Error
	if err != nil {
		// 角色增加失败
		con.Error(c, "角色增加失败", "/admin/role/add")
	} else {
		con.Success(c, "角色增加成功", "/admin/role")
	}
}
func (con RoleController) Edit(c *gin.Context) {
	// Query查询?后的数据
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/role")
	} else {
		role := models.Role{Id: id}
		models.DB.Find(&role)
		c.HTML(http.StatusOK, "admin/role/edit.html", gin.H{
			"role": role,
		})
	}

}
func (con RoleController) DoEdit(c *gin.Context) {
	//查询要修改的数据
	id, err1 := models.Int(c.PostForm("id"))
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/role")
	} else {
		role := models.Role{Id: id}
		models.DB.Find(&role)
		role.Title = title
		role.Description = description
		err2 := models.DB.Save(&role).Error
		if err2 != nil {
			con.Error(c, "修改角色信息失败", "/admin/role/edit?id="+models.String(id))
		} else {
			con.Success(c, "修改角色信息成功", "/admin/role")
		}
	}

}
func (con RoleController) Delete(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/role")
	} else {
		role := models.Role{Id: id}
		models.DB.Delete(&role)
		con.Success(c, "删除角色信息成功", "/admin/role")
	}
}

func (con RoleController) Auth(c *gin.Context) {
	// 获取当前角色Id
	roleId, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/role")
		return
	}
	// 获取所有权限
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem").Find(&accessList)

	// 查询当前角色拥有的权限，把权限id放到一个map对象里
	roleAccessList := []models.RoleAccess{}
	models.DB.Where("role_id=?", roleId).Find(&roleAccessList)
	roleAccessMap := make(map[int]int) //创建map类型对象
	for _, v := range roleAccessList {
		roleAccessMap[v.AccessId] = v.AccessId
	}
	// 遍历所有权限Id，判断是否在当前map对象里，是的话加上一个checked属性
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[accessList[i].Id]; ok {
			accessList[i].Checked = true
		}
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	// 注意静态页面不要/admin,而是admin
	c.HTML(200, "admin/role/auth.html", gin.H{
		"roleId":     roleId,
		"accessList": accessList,
	})
}
func (con RoleController) DoAuth(c *gin.Context) {
	// 获取角色id
	roleId, err1 := models.Int(c.PostForm("roleId"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/role")
		return
	}
	// 获取权限id
	accessIds := c.PostFormArray("access_node[]")

	// 删除当前角色的所有权限
	roleAccess := models.RoleAccess{}
	models.DB.Where("role_id=?", roleId).Delete(&roleAccess)

	// 遍历权限id,放入role_access表中
	for _, v := range accessIds {
		roleAccess.RoleId = roleId
		accessId, err2 := models.Int(v)
		if err2 != nil {
			con.Error(c, "传入参数错误", "/admin/role")
			return
		}
		roleAccess.AccessId = accessId
		models.DB.Create(&roleAccess)
	}
	fmt.Println(roleId)
	fmt.Println(accessIds)
	c.String(200, "doauth")
}
