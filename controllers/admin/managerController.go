package admin

import (
	"go_demo/models"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {
	// 这里需要关联查询才能显示角色信息
	managerList := []models.Manager{}
	// Role是manager结构体中的Role
	models.DB.Preload("Role").Find(&managerList)
	//fmt.Printf("%#v", managerList)
	c.HTML(200, "admin/manager/index.html", gin.H{
		"managerList": managerList,
	})
}
func (con ManagerController) Add(c *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(200, "admin/manager/add.html", gin.H{
		"roleList": roleList,
	})
}
func (con ManagerController) Edit(c *gin.Context) {
	// 首先获取管理员id(Url传值)
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "获取管理员id失败", "/admin/manager")
		return
	}
	// 获取所有角色
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	// 获取当前管理员
	manager := models.Manager{Id: id}
	models.DB.Find(&manager)
	// MD5解密
	manager.Password=models.MD5(manager.Password)
	c.HTML(200, "admin/manager/edit.html", gin.H{
		"roleList": roleList,
		"manager":  manager,
	})
}
func (con ManagerController) DoEdit(c *gin.Context) {
	managerId, err1 := models.Int(c.PostForm("id"))
	// 接收表单managerId
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/manager")
		// 防止继续往下执行
		return
	}
	roleId,err2:=models.Int(c.PostForm("role_id"))
	// 接收表单roleId
	if err2!=nil{
		con.Error(c, "传入参数错误", "/admin/manager")
		// 防止继续往下执行
		return
	}
	// 接收表单其他信息
	username := c.PostForm("username")
	password := c.PostForm("password")
	mobile := c.PostForm("mobile")
	email := c.PostForm("email")
	
    // 修改并保存
	manager := models.Manager{Id: managerId}
	models.DB.Find(&manager)
	manager.Username = username
	if password!=""{
		// 判断密码长度是否合法
		if len(password)<6{
			// 注意managerId是int类型，不能写成"/admin/manager/edit?id=managerId"
			con.Error(c,"密码长度不能小于6","/admin/manager/edit?id="+models.String(managerId))
			return
		}
		manager.Password = models.MD5(password)
	}
	manager.Email = email
	manager.Mobile = mobile
	manager.RoleId = roleId
	err3 := models.DB.Save(&manager).Error
	if err3 != nil {
		con.Error(c, "修改管理员信息失败", "/admin/manager/edit?id="+models.String(managerId))
	} else {
		con.Success(c, "修改管理员信息成功", "/admin/manager")
	}
}
func (con ManagerController) Delete(c *gin.Context) {
	managerId,err1:=models.Int(c.Query("id"))
	if err1!=nil{
		con.Error(c,"传入参数错误","/admin/manager")
	}
	manager:=models.Manager{Id: managerId}
	err2:=models.DB.Delete(&manager).Error
	if err2!=nil{
		con.Error(c,"删除管理员失败","/admin/manager")
	}else{
		con.Success(c,"删除管理员成功","/admin/manager")
	}
}
func (con ManagerController) DoAdd(c *gin.Context) {
	// 注意表单传过来的是string类型的
	roleId, err1 := models.Int(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/manager/add")
		// 防止继续往下执行
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	mobile := c.PostForm("mobile")
	email := c.PostForm("email")
	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或者密码长度不合法", "/admin/manager/add")
		return
	}

	// 判断管理员是否存在
	managerList := []models.Manager{}
	models.DB.Where("username=?", username).Find(&managerList)
	if len(managerList) > 0 {
		con.Error(c, "管理员已经存在", "/admin/manager/add")
		return
	}
	manager := models.Manager{
		Username: username,
		Password: models.MD5(password),
		Email:    email,
		Mobile:   mobile,
		RoleId:   roleId,
		AddTime:  int(models.GetUnix()),
	}

	err2 := models.DB.Create(&manager).Error
	if err2 != nil {
		con.Error(c, "数据上传失败", "/admin/manager/add")
		return
	}
	con.Success(c, "增加管理员成功", "/admin/manager")

}
