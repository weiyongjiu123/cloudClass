package teacherController

import (
	"cloudClass/controllers"
	"cloudClass/models"
	"cloudClass/service"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

type TeacherController struct {
	controllers.BaseController
	teacherService service.TeacherService
	isHaveReturn bool
}



func (this *TeacherController) ControllerInit() {
	this.teacherService = service.TeacherService{}
}





func (this *TeacherController) Login() {
	body := this.Ctx.Input.RequestBody
	data := make(map[string]string)
	err := json.Unmarshal(body,&data)

	if err != nil{
		this.Data["json"] = map[string]int8{
			"code":0,
		}
		this.ServeJSON()
		return
	}

	teacherService := service.TeacherService{}
	loginRes ,teacherMap := teacherService.Login(data["username"], data["password"])

	if(loginRes){
		this.SetSession("login",true)
		this.SetSession("teacher",teacherMap)
		this.Data["json"] = map[string]int8{
			"code":1,
		}
	}else{
		this.Data["json"] = map[string]int8{
			"code":0,
		}
	}
	this.ServeJSON()
}

func (this *TeacherController)Session() {
	teacher := this.GetSession("teacher")
	login  := this.GetSession("login")
	registerValidTime := this.GetSession("registerValidTime")
	registerValidCode := this.GetSession("registerValidCode")

	fmt.Println(teacher,login,registerValidTime,registerValidCode)
	s := make(map[string]interface{})
	s["login"] = login
	s["teacher"] = teacher
	s["registerValidTime"] = registerValidTime
	s["registerValidCode"] = registerValidCode
	this.Data["json"] = s
	this.ServeJSON()
}

func (this *TeacherController) returnJson(status int, content interface{})  {
	if content == nil{
		content = ""
	}
	res := make(map[string]interface{})
	res["code"] = status
	res["content"] = content
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *TeacherController) returnJsonWithTeacher(status int, content interface{}) {
	teacher := this.GetSession("teacher").(models.Teacher)
	data := map[string]interface{}{
		"code":status,
		"teacher": map[string]interface{}{
			"amount":teacher.Amount,
			"img":teacher.Img,
			"name":teacher.Name,
			"phone":teacher.Phone,
			"createTime":teacher.CreateTime,
		},
		"content":content,
	}
	this.Data["json"] = data
	this.ServeJSON()
}


func (this *TeacherController) Register() {
	body := this.Ctx.Input.RequestBody
	data := make(map[string]string)
	err := json.Unmarshal(body,&data)
	if err != nil || len(data["phone"]) == 0 || len(data["password"]) == 0{
		this.returnJson(0,"数据错误")
		return
	}
	//验证验证码
	registerValidTime := this.GetSession("registerValidTime")
	registerValidCode := this.GetSession("registerValidCode")
	if registerValidCode != data["code"]{
		this.returnJson(0,"验证码不正确")
		return
	}
	if(registerValidTime.(int64) > utils.GetTime()){
		this.returnJson(0,"验证码超时")
		return
	}


	//注册教师信息
	teacherService := service.TeacherService{}
	isHaveUser := teacherService.IsHaveUser(data["phone"])
	if(isHaveUser){
		this.returnJson(0,"该手机已经注册过")
		return
	}
	registerRes,teacherMap := teacherService.Register(data["phone"], data["password"])
	if(registerRes){
		this.SetSession("teacher",teacherMap)
		this.Data["json"] = map[string]int8{
			"code":1,
		}
	}else{
		this.Data["json"] = map[string]int8{
			"code":0,
		}
	}
	this.ServeJSON()

}

type res struct {
	Code int
	Content string
}

func (this *TeacherController) GetCode() {
	phone := this.GetString("phone")
	code ,sendCodeRes:= utils.SendCode(phone)
	 codeRes := 0
	if(sendCodeRes){
		codeRes = 1
		time, _ := beego.AppConfig.Int64("appRegisterValidTime")
		this.SetSession("registerValidTime",utils.GetTime() + time)
		this.SetSession("registerValidCode",code)

	}
	re := res{codeRes,code}
	this.Data["json"] = re
	this.ServeJSON()

}

func (this * TeacherController)DelSe()  {
	//this.DestroySession()
}






