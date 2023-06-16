package service

import (
	"IM/models"
	"IM/sql"
	"github.com/gin-gonic/gin"
)

type GroupService struct {
	sql.Groupdao
}

func (T *GroupService) Router(engine *gin.Engine) {
	group := engine.Group("/Group")
	group.GET("test", T.test)
	group.POST("addgroup", T.createGroup)
	group.POST("delGroup", T.delGroup)
}

func (T *GroupService) test(context *gin.Context) {
	context.JSON(200, models.GroupBasic{})

}

func (T *GroupService) createGroup(context *gin.Context) {
	group := models.GroupBasic{}
	context.BindJSON(&group)
	b := T.Insert(group)
	context.JSON(200, &b)
}

func (T *GroupService) delGroup(context *gin.Context) {

}
