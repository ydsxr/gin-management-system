package main

// 单引号字符，双引号字符串
import (
	"go_demo/models"
	"go_demo/routers"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Article struct {
	Title   string
	Desc    string
	Content string
}

// 自定义函数
func UnixToTime(timestamp int) string {
	//传入毫秒、微秒
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func main() {
	// 创建了一个 gin.Engine 实例（核心结构体，负责整个web服务的调度）
	r := gin.Default()
	// 静态文件路径
	r.Static("/static", "./static")
	store := cookie.NewStore([]byte("secret111"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600, // Cookie 有效期（秒）
		HttpOnly: true,
		// Domain: "10.100.233.37", //  指定域名
	})
	r.Use(sessions.Sessions("mysession", store))

	//自定义模板函数 html中用到的函数
	r.SetFuncMap(template.FuncMap{
		"UnixToTime": UnixToTime,
		"Str2Html":   models.Str2Html,
	})
	r.LoadHTMLGlob("templates/**/**/*")

	routers.AdminRoutersInit(r)
	//routers.ApiRoutersInit(r)
	routers.DefaultRoutersInit(r)
	r.Run(":8080") //启动web服务
}
