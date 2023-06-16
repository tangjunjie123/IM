/**
* @Auth:ShenZ
* @Description:
* @CreateDate:2022/06/15 14:57:55
 */
package models

import (
	"gorm.io/gorm"
)

//人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁 /群 ID
	Type     int  //对应的类型  1好友  2群  3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
