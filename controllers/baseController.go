package controllers

import (
	"cloudClass/public"
	"cloudClass/request"
	"cloudClass/service"
	"cloudClass/utils"
	"fmt"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	BaseControllerName string	//控制器名称 TeacherController
	BaseActionName string		//方法名	GET:Login
	BaseRequestType string		//请求方式 只有两种，ordinary 和 permissions
	BaseReqData request.ReqData	//post 数据
	BaseService service.BaseService
	BaseCtrlMethod string
}


//var validOpen, _ = beego.AppConfig.Bool("validopen")
var validOpen = true

func (this *BaseController) baseInit() {
	//controller, action := this.GetControllerAndAction()
	arr,ok := utils.StrSeg(this.Ctx.Input.URL(),"^/([\\S\\s]*)/([\\S\\s]*)$")
	if !ok{
		this.ReturnContent(-3,"系统繁忙，请重试")
	}
	//fmt.Println(this.Ctx.Input.URL())
	//method := utils.StrMerge(this.Ctx.Request.Method, ":")
	//action = utils.StrMerge(method, action)
	this.BaseControllerName = arr[1]
	this.BaseActionName = this.Ctx.Input.URL();
}


func (this *BaseController) initSession() {
	//if (validOpen) {

		fmt.Println("init session .....")
		//this.SetSession("login", true)
		//teacher := models.Teacher{100, "小花", "12345678910", "123456", 0, 123456, "/static/img/teacher/default.jpg",0,0,0,0}
		//oeg225I9auRLzcI4QMXELJdVYKJc   testOpenId
		//user := models.User{"oeg225I9auRLzcI4QMXELJdVYKJc","王小明",100,132,"王小明","13213",100}
		//this.SetSession("user","test")
		//this.SetSession("teacher", "test")
		//this.SetSession("admin","admin")

	//	validOpen = false;
	//}
}

func (this *BaseController) validPermissions() {

	RouterMap := public.RouterMap[this.BaseControllerName]
	moduleMap := RouterMap["permissions"].(map[string]interface{})
	this.BaseRequestType = "ordinary"

	for k, _ := range moduleMap {
		//method := (v.(map[string]interface{}))["method"].(string)
		if k == this.BaseActionName{
			this.BaseRequestType = "permissions"
			if user := this.GetSession(public.ModuleSessionMap[this.BaseControllerName]["name"]); user == nil {
				this.Data["json"] = map[string]interface{}{
					"code":    -2,
					"content": "权限不够",
					"role":this.BaseControllerName,
				}
				this.ServeJSON()
				this.StopRun()
			}
		}
	}
}

func (this *BaseController) getPostData() {
	//只有post this.BaseCtrlMethod = "reqDataValid"
	if this.Ctx.Request.Method != "POST" || this.BaseCtrlMethod != "reqDataValid" {
		return
	}
	//fmt.Println(this.BaseControllerName,this.BaseRequestType,this.Ctx.Request.RequestURI)
	if len(this.Ctx.Input.RequestBody)!=0{
		parseRes := this.BaseReqData.ParseData(this.Ctx.Input.RequestBody)
		if !parseRes{
			this.baseStopRun(0,"参数错误")
		}
	}

	res := this.BaseReqData.Valid()
	if !res {
		this.baseStopRun(0,"参数错误")
	}
}

func (this *BaseController) ControllerInit()  {

}

func (this *BaseController) ReturnWithPerSession(status int8,content interface{}) {
	//fmt.Println("---teacher : ",this.GetSession("teacher"))
	name := public.ModuleSessionMap[this.BaseControllerName]["name"].(string)
	returnSession := public.ModuleSessionMap[this.BaseControllerName]["returnSession"].(func(n string,i interface{})map[string]interface{})
	session := returnSession(name,this.GetSession(name))

	data := map[string]interface{}{
		"code":status,
		name: session,
		"content":content,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) ReturnContent(status int8, content interface{})  {
	res := make(map[string]interface{})
	res["code"] = status
	res["content"] = content
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *BaseController)baseStopRun(status int,content interface{})  {
	this.Data["json"] = map[string]interface{}{
		"code":    status,
		"content": content,
	}
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) reqDataInit() {

}

func (this *BaseController) setCtrlMethod() {
	moduleMap := public.RouterMap[this.BaseControllerName]
	controllerMap := moduleMap[this.BaseRequestType].(map[string]interface{})
	urlArr,ok := utils.StrSeg(this.Ctx.Request.RequestURI,"^([\\S\\s]*)\\?")
	url := this.Ctx.Request.RequestURI
	//路径是否携带参数
	if ok{
		url = urlArr[1]
	}
	//fmt.Println(this.BaseControllerName,controllerMap,url)
	ctrlAttrMap := controllerMap[url].(map[string]interface{})
	fun,ok :=ctrlAttrMap["reqDataValid"]
	if ok{
		this.BaseCtrlMethod = "reqDataValid"
		this.BaseReqData = (fun.(func()interface{})()).(request.ReqData)////
	}
}


func (this *BaseController) Prepare() {
	//应用初始化
	this.baseInit()
	this.initSession()
	this.validPermissions()
	this.setCtrlMethod()
	//ctrl第二方式初始化
	this.BaseReqDataInit()
	this.getPostData()
	//ctrl第一方式
	this.ControllerInit()

}

func (this *BaseController) BaseReqDataInit() {
	if this.BaseCtrlMethod != "reqDataValid" {
		return
	}
	ctrl := new(request.Ctrl)
	ctrl.GetSessionFun = func(name string) interface{} {
		return this.GetSession(name)
	}
	ctrl.GetInt64 = this.GetInt64
	ctrl.GetParIntFun = func(name string) int {
		value,_ := this.GetInt(name)
		return value
	}
	ctrl.GetParStringFun = func(name string) string {
		return this.GetString(name)
	}
	ctrl.SetSessionFun = this.SetSession
	this.BaseReqData.Init(ctrl,this.BaseReqData,this.BaseRequestType)

	ctrl.SaveFile = func(url string, key string) error {
		f, _, err := this.GetFile(key)
		if err != nil {
			return err
		}
		defer f.Close()
		return this.SaveToFile(key, url)
	}
	ctrl.Redirect = func(url string, code int) {
		this.Redirect(url,code)
	}
}



