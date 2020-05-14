package request

import (
	"cloudClass/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type ReqData interface {
	Valid()(res bool)
	validArr()([]int,[]int64,[]string)
	ParseData(bytes []byte)bool
	Prepare()
	Model()
	Service()
	ReturnRes() (int8,interface{})
	Check() bool
	Init(ctrl *Ctrl,reqData ReqData,t string)
	End()
	GetFun() []func()
	Next(_if func(bool), _true func([]func()),_false func([]func()))
	Clear()
	IsPlay() bool
	GetVideoUrl() string
}

/*------------------------ 模板----------------------------*/
type BaseReq struct {
	baseStatus int8
	baseContent interface{}
	isReturn bool
	ctrl *Ctrl
	reqData ReqData
	checkErr error
	checkBool bool
	checkZero int64
	checkString string
	isPlayVideo bool
	videoUrl string
	check bool
}

func (this *BaseReq)Valid()(res bool) {
	ints, int64s, strings := this.reqData.validArr()
	for e := range ints {
		if ints[e] == 0{
			return false
		}
	}
	for e := range int64s {
		if int64s[e] == 0{
			return false
		}
	}
	for e := range strings {
		if len(strings[e]) == 0{
			return false
		}
	}
	return true
}
func (this *BaseReq)ParseData(bytes []byte)bool {
	err := json.Unmarshal(bytes,this.reqData)
	if err != nil{
		return false
	}
	return true
}
func (this *BaseReq) Prepare(){}

func (this *BaseReq)Model()  {}

func (this *BaseReq) Service ()(){}
func (this *BaseReq)setReturn(status int8 ,content interface{})  {
	this.baseStatus = status
	this.baseContent = content
	this.isReturn = true
}

func(this *BaseReq) Clear() ()  {
	this.baseContent = ""
	this.baseStatus = -100
	this.check = true
}

func (this *BaseReq)IsPlay() bool {
	return this.isPlayVideo
}

func (this *BaseReq) ReturnRes() (int8,interface{}){
	return this.baseStatus,this.baseContent
}

func (this *BaseReq) isCheck()bool  {
	if !this.check {
		this.checkZero = 1
		this.checkErr = nil
		this.checkBool = true
		this.checkString = "123"
	}
	return this.check
}

func (this *BaseReq) Check() bool{
	if !this.isCheck(){
		return this.isReturn
	}
	if this.checkErr != nil{
		fmt.Println(this.checkErr)
	}
	if this.checkZero == 0 || this.checkErr != nil || !this.checkBool  || len(this.checkString) == 0{
		if this.baseStatus == -100{
			this.setReturn(0,this.baseContent)
		}else{
			this.setReturn(this.baseStatus,this.baseContent)
		}
		return true
	}
	return this.isReturn
}

func (this *BaseReq)checkInit()  {

	this.checkErr = nil
	this.checkBool = true
	this.checkZero = 1
	this.checkString = "我不是空的"
	this.baseContent = ""
	this.baseStatus = -100
	this.check = true
}

func (this *BaseReq) validInt64Arr(arr []int64)  {
	for e := range arr {
		if arr[e] == 0 {
			this.checkZero = 0
		}
	}
}

func (this *BaseReq) Init(ctrl *Ctrl,reqData ReqData) {
	this.ctrl = ctrl
	this.reqData = reqData
	this.checkInit()
}
func (this *BaseReq) GetFun() []func() {
	return []func(){
		this.reqData.Prepare,
		this.reqData.Model,
		this.reqData.Service,
		this.reqData.End,
	}
}

func (this *BaseReq) setVideoUrl(url string) {
	this.isPlayVideo = true
	this.videoUrl = url
	this.isReturn = true
}

func (this *BaseReq) GetVideoUrl() string {
	return this.videoUrl
}

//func (this *BaseReq) conversion(v interface{},t string) {
//	switch t {
//	case "int":
//
//	}
//}
//如何前面都没有回复，默认是操作失败的
func (this *BaseReq)End()   {

}

func (this *BaseReq) resConver(res bool)  int8{
	if res{
		return 1
	}else{
		return 0
	}
}

func (this *BaseReq)validArr()([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{}
}

func (this *BaseReq) setReturnBool(b bool,content interface{}) {
	this.setReturn(this.resConver(b),content)
}

func (this *BaseReq) setReturnNoContent (b bool){
	this.setReturn(this.resConver(b),"")
}

func (this *BaseReq) setReturnWithCount (b bool){
	this.setReturn(this.resConver(b),"操作失败，请重试")
}

func (this *BaseReq)sqlQuery(sql string,arr interface{},params...interface{}) (int64,error) {
	o := orm.NewOrm()
	return o.Raw(sql,params).QueryRows(arr)
}

func (this *BaseReq)judgeEq(v interface{}, v2 interface{}) {
	if v == v2{
		this.checkZero = 0
	}
}

func (this *BaseReq)judge(b bool)  {
	if b{
		this.checkZero = 0
	}
}

func (this *BaseReq)judgeNeq(v interface{}, v2 interface{}) {
	if v != v2{
		this.checkZero = 0
	}
}

func (this *BaseReq)isInArr(v interface{},arr[]interface{}) {
	for e := range arr {
		if arr[e] == v{
			this.checkZero = 0
		}
	}
}

func (this *BaseReq)isNotInArr(v interface{},arr[] interface{}) {
	this.checkZero = 0
	for e := range arr {
		//fmt.Println(arr[e] == v,arr[e],v,utils.Typeof(arr[e]),utils.Typeof(v))
		if arr[e] == v{
			this.checkZero = 1
		}
	}

}
func (this *BaseReq)Next(_if func(bool), _true func([]func()),_false func([]func())){

}


/*----------------------------------------------------*/

type TeacherBase struct {
	BaseReq
	Teacher models.Teacher
}

func (this *TeacherBase) Init(ctrl *Ctrl,reqData ReqData,baseRequestType string) {
	this.checkInit()
	this.ctrl = ctrl
	this.reqData = reqData
	if baseRequestType == "permissions" {
		this.Teacher = this.ctrl.GetSession("teacher").(models.Teacher)
	}
	//var teacherList []models.Teacher
	//this.checkZero,this.checkErr = this.sqlQuery("select * from teacher where id = 100",&teacherList)
	//this.Teacher = teacherList[0]
	//this.ctrl.SetSession("teacher",teacherList[0])
}

