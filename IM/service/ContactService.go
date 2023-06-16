package service

import (
	"IM/models"
	"IM/sql"
	"github.com/gin-gonic/gin"
)

type ContactService struct {
	sql.Contactdao
}

func (T *ContactService) Router(gin *gin.Engine) {
	group := gin.Group("/contact")
	group.POST("addcontact", T.AddFrind)
	group.GET("test", T.Test)
	group.POST("delcontact", T.DelContact)
}
func (T *ContactService) Test(gin *gin.Context) {
	gin.JSON(200, models.Contact{})
}

func (T *ContactService) AddFrind(gin *gin.Context) {
	contact := models.Contact{}
	gin.BindJSON(&contact)
	b := T.Insert(contact)
	gin.JSON(200, &b)
}

func (T *ContactService) DelContact(gin *gin.Context) {
	contact := models.Contact{}
	gin.BindJSON(&contact)
	T.Del(contact)
}
