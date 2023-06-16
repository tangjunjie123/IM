package sql

import (
	"IM/models"
	"IM/utils"
	"fmt"
)

type Contactdao struct {
}

func (T *Contactdao) Insert(contact models.Contact) bool {
	db := utils.Db
	tx := db.Create(&contact)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return false
	}
	return true
}

func (T *Contactdao) Find2(id int64) []models.Contact {
	ans := []models.Contact{}
	db := utils.Db
	db.Where("Owner_Id = ? && Type = 2 ", id).Find(&ans)
	return ans
}

func (T *Contactdao) Find1(id int) []models.Contact {
	ans := []models.Contact{}
	db := utils.Db
	db.Where("Owner_Id = ? && Type = 1", id).Find(&ans)
	return ans
}

func (T *Contactdao) Del(contact models.Contact) {

	db := utils.Db

	db.Where("Owner_Id = ?  && Target_Id =? && Type =?", contact.OwnerId, contact.TargetId, contact.Type).Delete(&contact)
	//db.Delete(&contact)
}
