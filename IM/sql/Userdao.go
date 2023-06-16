package sql

import (
	"IM/models"
	"IM/utils"
)

type Userdao struct {
}


func (Userdao) Find(user models.UserBase) models.UserBase {
	db := utils.Db
	ans := models.UserBase{}
	db.Where("Name", user.Name).Find(&ans)
	return ans
}

func (Userdao) Insert(user ...models.UserBase) int {
	db := utils.Db
	s := 0
	for _, u := range user {
		db.Create(&u)
		s++
	}
	return s
}

func (Userdao) FindById(id int64) models.UserBase {
	db := utils.Db
	u := models.UserBase{}
	db.Find(&u, id)
	return u
}
