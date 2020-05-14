package main

import (
	"cloudClass/models"
	_ "cloudClass/routers"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	models.Init()
}
func main() {

	beego.Run()
}

