package main

import (
	"auth-go-example/middlewares"
	"auth-go-example/models"
	_ "auth-go-example/routers"
	"github.com/astaxie/beego"
)

func main() {
	models.MigrationDB()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
	}

	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.InsertFilter("/*", beego.BeforeRouter, middlewares.AuthMiddleware)
	beego.Run()
}
