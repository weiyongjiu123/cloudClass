package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

/*------------------------ 模板 A----------------------------*/

type BaseReqA struct {
	BaseReq
}

func (this *BaseReqA)sqlQuery(sql string,arr interface{},params...interface{})   {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	o := orm.NewOrm()
	this.checkZero,this.checkErr =  o.Raw(sql,params).QueryRows(arr)
}


func (this *BaseReqA)sqlIsExist(sql string,params...interface{})   {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	var maps []orm.Params
	o :=orm.NewOrm()
	var num int64
	num ,this.checkErr = o.Raw(sql,params).Values(&maps)
	this.judge(num>0)
}

func (this *BaseReqA)sqlTransaction(funArr ...func(o orm.Ormer, res *bool)) {
	this.checkBool = utils.SqlTransaction(funArr...)
}

func (this *BaseReqA)add(table string, fields map[string]interface{})int64 {
	fieldStr := ""
	mapLen := len(fields)
	valuesStr := ""
	var valueArr []interface{}
	var lastId int64

	count := 0
	for k, v := range fields {
		count++
		valueArr = append(valueArr,v)
		if count == mapLen{
			fieldStr += k
			valuesStr += "?"
		}else{
			fieldStr += k + ","
			valuesStr += "?,"
		}
	}
	sql := "insert into " + table + "("+fieldStr+") values("+valuesStr+") "
	this.checkBool, lastId, _ = utils.SqlExec(sql, valueArr)
	return lastId
}

func (this *BaseReqA)update(table string, fields map[string]interface{},whereArr map[string]interface{},condition string)  {
	updateFields := ""
	whereStr := ""
	mapLen := len(fields)
	whereLent := len(whereArr)

	var valueArr []interface{}

	count := 0
	for k, v := range fields {
		count++
		valueArr = append(valueArr,v)
		if count == mapLen{
			updateFields += k + "=?"
		}else{
			updateFields += k + "=?,"
		}
	}
	count = 0
	for k, v := range whereArr {
		count++
		valueArr = append(valueArr,v)
		if count == whereLent{
			whereStr += k + "=?"
		}else{
			whereStr += k + "=? " + condition +" "
		}
	}

	sql := "update " + table + " set " + updateFields + " where "+whereStr
	this.checkBool, _, _ = utils.SqlExec(sql, valueArr)

}

func (this *BaseReqA)del(table string,whereArr map[string]interface{},condition string)  {
	whereStr := ""
	whereLent := len(whereArr)
	count := 0
	var valueArr []interface{}

	for k, v := range whereArr {
		count++
		valueArr = append(valueArr,v)
		if count == whereLent{
			whereStr += k + "=?"
		}else{
			whereStr += k + "=? " + condition +" "
		}
	}
	sql := "delete from " + table + " where " +whereStr
	this.checkBool, _,_  = utils.SqlExec(sql, valueArr)

}

func (this *BaseReqA) sqlExec(sql string,params... interface{})int64  {
	var lastId int64
	res,lastId,_ := utils.SqlExec(sql,params...)
	this.checkBool = res
	return lastId
}

func (this *BaseReqA)errMsg(msg string)  {
	this.baseStatus = -1
	this.baseContent = msg
}


/*------------------------ 模板 A----------------------------*/


type TeacherBaseA struct {
	BaseReqA
	Teacher models.Teacher
}

func (this *TeacherBaseA) Init(ctrl *Ctrl,reqData ReqData,baseRequestType string) {
	this.checkInit()
	this.ctrl = ctrl
	this.reqData = reqData
	if baseRequestType == "permissions" {
		this.Teacher = this.ctrl.GetSession("teacher").(models.Teacher)
	}
	//var teacherList []models.Teacher
	//this.sqlQuery("select * from teacher where id = 100",&teacherList)
	//this.Teacher = teacherList[0]
	//this.ctrl.SetSession("teacher",teacherList[0])
}

func (this *TeacherBaseA) getTeacherId()string {
	return strconv.FormatInt(this.Teacher.Id,10)
}

func (this *TeacherBaseA)flushTeacher() {
	var teacherList  []models.Teacher
	this.sqlQuery("select * from `teacher` where id = ?",&teacherList,this.Teacher.Id)
	this.Teacher = teacherList[0]
	this.ctrl.SetSession("teacher",this.Teacher)
}



/*-----------------------------------------------------------------*/
type UserA struct {
	BaseReqA
	user models.User
}

func (this *UserA) Init(ctrl *Ctrl,reqData ReqData,baseRequestType string) {
	this.checkInit()
	this.ctrl = ctrl
	this.reqData = reqData
	if baseRequestType == "permissions" {
		this.user = this.ctrl.GetSession("user").(models.User)
	}
	//var userList []models.User
	//oeg225I9auRLzcI4QMXELJdVYKJc testOpenId
	//this.sqlQuery("select * from `user` where id=?",&userList,"oeg225I9auRLzcI4QMXELJdVYKJc")
	//this.user = userList[0]
}

func (this *UserA)flushUser() {
	var userList  []models.User
	this.sqlQuery("select * from `user` where id = ?",&userList,this.user.Id)
	this.user = userList[0]
	this.ctrl.SetSession("user",this.user)
}

/*---------------------------------------------------*/
type AdminA struct {
	BaseReqA
	Admin models.Admin
}

func (this *AdminA) Init(ctrl *Ctrl,reqData ReqData,baseRequestType string) {
	this.checkInit()
	this.ctrl = ctrl
	this.reqData = reqData
	if baseRequestType == "permissions" {
		this.Admin = this.ctrl.GetSession("admin").(models.Admin)
	}
	//var adminList []models.Admin
	//this.sqlQuery("select * from `admin` where id = 1",&adminList);
	//this.Admin = adminList[0];
}



