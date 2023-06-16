package main

import (
	"IM/models"
	"fmt"
	"gorm.io/driver/mysql" // gorm mysql 驱动包
	"gorm.io/gorm"         // gorm
)

func main() {
	username := "root"     // 账号
	password := "Tjj@2002" // 密码
	host := "123.56.9.154" // 地址
	port := 3306           // 端口
	DBname := "IM"         // 数据库名称
	timeout := "10s"       // 连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
	// Open 连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql.")
	}
	fmt.Println(db)
	db.AutoMigrate(&models.GroupBasic{})
}
