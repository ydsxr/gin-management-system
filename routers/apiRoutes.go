package routers

import "github.com/gin-gonic/gin"

func ApiRoutersInit(r *gin.Engine) {
	adminRouters := r.Group("/api")
	{
		adminRouters.GET("/", func(c *gin.Context) {
			c.String(200, "api首页")
		})
		adminRouters.GET("/user", func(c *gin.Context) {
			c.String(200, "api用户列表")
		})
		adminRouters.GET("/article", func(c *gin.Context) {
			c.String(200, "api新闻列表")
		})
	}
}