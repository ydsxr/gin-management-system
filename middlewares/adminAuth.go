package middlewares

import (
	"encoding/json"
	"fmt"
	"go_demo/models"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

func InitAdminAuthMiddleware(c *gin.Context) {
	//权限判断 没有登录的用户不能进入后台管理中心
	// excludeAuthPath("aaa")
	// 1,获取访问的URL地址
	// /admin/captcha?t=0.5362956341580839
	// pathname:=c.Request.URL.String()
	// fmt.Println(pathname)
	// 只要前面的部分 /admin/captcha
	pathname := strings.Split(c.Request.URL.String(), "?")[0]
	fmt.Println(pathname)

	// 2，获取session中保存的信息
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")
	// 类型断言 userinfo是不是string类型
	userinfoStr, ok := userinfo.(string)
	if ok {
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		// 如果没有登录成功
		if !(len(userinfoStruct) > 0 && userinfoStruct[0].Username != "") {
			if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
				// 用户没有登录，跳转到login页面
				c.Redirect(302, "/admin/login")
			}
		} else { //登陆成功后的权限判断
			urlPath := strings.Replace(pathname, "/admin/", "", 1)
			// 判断是否是超级管理员和通用URL
			if userinfoStruct[0].IsSuper == 0 && !excludeAuthPath("/"+urlPath){
				//1、根据角色获取当前角色的权限列表，把权限id放在map类型的对象里
				roleAccessList := []models.RoleAccess{}
				models.DB.Where("role_id=?", userinfoStruct[0].RoleId).Find(&roleAccessList)
				roleAccessMap := make(map[int]int)
				for _, v := range roleAccessList {
					roleAccessMap[v.AccessId] = v.AccessId
				}
				//2、根据当前访问的Url对应的权限id，判断该id是否在当前角色的权限列表里
				access := models.Access{}
				models.DB.Where("url=?", urlPath).Find(&access)
				if _, ok := roleAccessMap[access.Id]; !ok {
					c.String(200, "没有权限")
					c.Abort()
				}
			}

		}
		// fmt.Println(userinfoStruct)
		// c.JSON(http.StatusOK, gin.H{
		// 	"username": userinfoStruct[0].Username,
		// })
	} else { // 用户没有登录
		// 首先排除不需要权限判断的路由
		if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
			// 用户没有登录，跳转到login页面
			c.Redirect(302, "/admin/login")
		}
	}
}

// 需要排除权限判断的网址
func excludeAuthPath(urlPath string) bool {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	excludeAuthPath := config.Section("").Key("excludeAuthPath").String()
	excludeAuthPathSlice := strings.Split(excludeAuthPath, ",")
	// fmt.Println(excludeAuthPathSlice)
	for _, v := range excludeAuthPathSlice {
		if v == urlPath {
			return true
		}
	}
	return false
}
