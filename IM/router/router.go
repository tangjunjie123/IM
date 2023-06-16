package router

import (
	"IM/service"
	"IM/utils"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin := gin.Default()
	utils.Mysql_init()
	utils.Redis_init()
	Routers(gin)
	return gin
}
func Routers(gin *gin.Engine) {
	new(service.UserService).Router(gin)
	new(service.ContactService).Router(gin)
	new(service.GroupService).Router(gin)
}
