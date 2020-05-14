package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

func Init() {
	//host := beego.AppConfig.String("dbhost")66
	//port := beego.AppConfig.String("dbport")
	//user := beego.AppConfig.String("dbuser")
	//password := beego.AppConfig.String("dbpassword")
	//name := beego.AppConfig.String("dbname")
	//if port == ""{
	//	port = "3306"
	//}
	//dsn := user + ":" + password + "@tcp("+host+":"+port+")/" + name + "?charset=utf8"
	dsn := "root:root@tcp(127.0.0.1:3306)/cloudclass?charset=utf8"
	fmt.Println(dsn)
	orm.RegisterDataBase("default","mysql",dsn)
}
