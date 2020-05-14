package routers

import (
	"cloudClass/controllers/teacherController"
	"cloudClass/public"
	"github.com/astaxie/beego"
)

//var RouterMap = map[string]string{
//	"/teacher/login":"post:Login",
//	"/teacher/testSession":"get:Session",
//	"/teacher/register":"post:Register",
//	"/teacher/getCode":"get:GetCode",
//	"/teacher/addClass":"post:AddClass",
//	"/teacher/getTeacherMsg":"get:GetTeacherMsg",
//	"/teacher/addClassInit":"get:AddClassInit",
//}

func init() {

	for _, v := range public.RouterMap {
		permissionsMap,_ := v["permissions"].(map[string]interface{})
		ordinaryMap,_ := v["ordinary"].(map[string]interface{})
		for k, v := range permissionsMap {
			method := (v.(map[string]interface{}))["method"].(string)
			beego.Router(k,&teacherController.TeacherController{},method)
		}
		for k, v := range ordinaryMap {
			method := (v.(map[string]interface{}))["method"].(string)
			beego.Router(k,&teacherController.TeacherController{},method)
		}
	}


	//teacherRouter := public.TeacherRouter
	//teacherRouterPer := public.TeacherRouterPer
	//utils.MergeMap(&teacherRouter,teacherRouterPer)
	//for k, v := range teacherRouter {
	//	beego.Router(k,&teacherController.TeacherController{},v)
	//}
    //beego.Router("/", &controllers.MainController{})
    //beego.Router("/teacher/login",&teacherController.TeacherController{},"post:Login")
    //beego.Router("/teacher/testSession",&teacherController.TeacherController{},"get:Session")
    //beego.Router("/teacher/register",&teacherController.TeacherController{},"post:Register")
    //beego.Router("/teacher/getCode",&teacherController.TeacherController{},"get:GetCode")
    //beego.Router("/teacher/addClass",&teacherController.TeacherController{},"post:AddClass")
    //beego.Router("/teacher/getTeacherMsg",&teacherController.TeacherController{},"get:GetTeacherMsg")



    //beego.Router("/teacher/index",&teacherController.TeacherController{},"get:IndexView")
    //beego.Router("/teacher/classManager",&teacherController.TeacherController{},"get:ClassManagerView")
	//beego.Router("/teacher/chapter",&teacherController.TeacherController{},"get:ChapterView")
	//beego.Router("/teacher/chapterAdd",&teacherController.TeacherController{},"get:ChapterAddView")
	//beego.Router("/teacher/updateVideo",&teacherController.TeacherController{},"post:UpdateVideo")
	//beego.Router("/teacher/test",&teacherController.TeacherController{},"get:TestView")
	//beego.Router("/teacher/addClass",&teacherController.TeacherController{},"get:AddClassView")
	//beego.Router("/teacher/chapterTitle",&teacherController.TeacherController{},"get:ChapterTitleView")
}
