package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {}

func (con BaseController) Success(c *gin.Context,message string, redirectUrl string) {
	c.HTML(http.StatusOK,"admin/public/success.html",gin.H{
		"message":message,
		"redirectUrl":redirectUrl,
	})
}
func (con BaseController) Error(c *gin.Context,message string, redirectUrl string) {
	c.HTML(http.StatusOK,"admin/public/error.html",gin.H{
		"message":message,
		"redirectUrl":redirectUrl,
	})
}
