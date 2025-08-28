package admin

import (
	"encoding/json"
	"fmt"
	"go_demo/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MainController struct {
}

func (con MainController) Index(c *gin.Context) {
	// 获取session
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")
	// 类型断言 userinfo是不是string类型
	userinfoStr, ok := userinfo.(string)
	if ok {
		//1、获取用户信息
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		//2、获取所有权限
		accessList := []models.Access{}
		// 降序DESC 升序ASC
		models.DB.Where("module_id=?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("access.sort ASC")
		}).Order("sort ASC").Find(&accessList)
		//3、获取当前用户的角色的权限
		roleAccessList := []models.RoleAccess{}
		models.DB.Where("role_id=?", userinfoStruct[0].RoleId).Find(&roleAccessList)
		//4、定义一个map对象，将当前角色的权限放进map对象
		roleAccessMap := make(map[int]int)
		for _, v := range roleAccessList {
			roleAccessMap[v.AccessId] = v.AccessId
		}
		//5、循环遍历
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

		fmt.Printf("%#v", accessList)
		fmt.Println(userinfoStruct[0].IsSuper)
		c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
			"username":   userinfoStruct[0].Username,
			"accessList": accessList,
			"isSuper":    userinfoStruct[0].IsSuper,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"": "session不存在",
		})
	}
}
func (con MainController) Welcome(c *gin.Context) {
	c.HTML(200, "admin/main/welcome.html", gin.H{})
}

// 公共修改状态的方法
func (con MainController) ChangeStatus(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入参数错误",
		})
		return
	}
	table := c.Query("table")
	field := c.Query("field")

	// Exec执行原生sql
	err2 := models.DB.Exec("update "+table+" set "+field+" = ABS("+field+"-1) where id=?", id).Error
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
	})
}
func (con MainController) ChangeNum(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入参数错误",
		})
		return
	}
	table := c.Query("table")
	field := c.Query("field")
	num := c.Query("num")

	err2 := models.DB.Exec("update "+table+" set "+field+" = "+num+" where id=?", id).Error
	if err2!=nil{
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改失败",
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "修改成功",
		})
	}

}
