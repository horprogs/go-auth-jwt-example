package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Database() *gorm.DB {
	db, err := gorm.Open("mysql", beego.AppConfig.String("dbuser")+":"+beego.AppConfig.String("dbpassword")+"@tcp(127.0.0.1:3306)/db?parseTime=true")
	if err != nil {
		panic(err)
	}

	return db
}
