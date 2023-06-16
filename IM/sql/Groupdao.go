package sql

import (
	"IM/models"
	"IM/utils"
	"fmt"
)

type Groupdao struct {
}

func (T *Groupdao) Insert(basic models.GroupBasic) bool {
	db := utils.Db
	b := models.GroupBasic{}
	db.Where("Owner_Id = ?", basic.OwnerId).Find(&b)

	if b.OwnerId != 0 {
		return false
	}
	tx := db.Create(&basic)
	if tx.Error != nil {
		fmt.Println(tx)
		return false
	}
	return true
}
func (T *Groupdao) Del() {

}
func (T *Groupdao) update() {

}
