package routers

import (
	"go_demo/controllers/itying"

	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/default")
	{
		defaultRouters.GET("/thumbnail1", itying.DefaultController{}.Thumbnail1)
		defaultRouters.GET("/goQrcode", itying.DefaultController{}.GoQrcode)
	}
}
