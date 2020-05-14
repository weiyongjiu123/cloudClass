package teacherController

import (
	"cloudClass/models"
	"cloudClass/request"
	"cloudClass/service"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (this *TeacherController) UpdateVideo(){
	f, _, err := this.GetFile("uploadname")
	if err != nil {
		log.Fatal("getfile err ", err)
	}
	defer f.Close()
	this.SaveToFile("uploadname", "static/upload/" + strconv.FormatInt(utils.GetTime(),10)+".avi") // 保存位置在 static/upload, 没有文件夹要先创建
}

func (this *TeacherController)getData(paraMap map[string]int)(dataMap map[string]interface{},e bool)  {
	body := this.Ctx.Input.RequestBody
	data := make(map[string]interface{})
	err := json.Unmarshal(body,&data)
	if err != nil{
		e = false
		return
	}
	for k, v := range paraMap {
		switch v {
		//整形不能为0
		case 1:
			if data[k] == 0{
				e = false
				return
			}
			//字符串不能为空
		case 2:
			if len(data[k].(string)) == 0{
				e = false
				return
			}
		}
	}
	return data,true
}

func (this *TeacherController) AddClass() {
	body := this.Ctx.Input.RequestBody
	data := request.AddClass{}
	err := json.Unmarshal(body,&data)
	if err != nil || !data.Valid(){
		this.returnJson(0,"请求参数出错")
		return
	}
	//保存图片
	time := utils.GetTime()
	randNum := utils.RandNum(100000, 999999)
	picName := utils.StrMerge(strconv.FormatInt(time,10),strconv.FormatInt(randNum,10))
	picUrl := utils.StrMerge(beego.AppConfig.String("classimg"),picName)
	picUrl = utils.StrMerge(picUrl,".jpg")
	savePicRes := utils.SaveBase64ToPic(data.ClassImg,picUrl)
	if !savePicRes{
		this.returnJson(0,"图片上传失败")
		return
	}
	data.ClassImg = utils.StrMerge("/",picUrl)
	var teacher models.Teacher
	teacherSession := this.GetSession("teacher")
	teacher = teacherSession.(models.Teacher)

	teacherService := service.TeacherService{}
	addRes := teacherService.AddClass(data,teacher.Id)
	//
	if addRes{
		this.returnJson(1,"")
	}else{
		this.returnJson(0,"")
	}
}

func (this *TeacherController) GetTeacherMsg() {
	teacher := this.GetSession("teacher").(models.Teacher)

	this.returnJson(1, map[string]interface{}{
		"teacher": map[string]interface{}{
			"amount":teacher.Amount,
			"createTime":teacher.CreateTime,
			"img":teacher.Img,
			"name":teacher.Name,
			"phone":teacher.Phone,
		},
	})

}

func (this *TeacherController) AddClassInit (){
	teacherService := service.TeacherService{}
	list,num := teacherService.GetClassType()
	if num > 0{
		this.returnJson(1, list)
	}else{
		this.returnJson(0,"")
	}
	return

}

func (this *TeacherController)videoPlay(url string)  {
	video, err := os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer video.Close()
	http.ServeContent(this.Ctx.ResponseWriter, this.Ctx.Request, "test.mp4", time.Now(), video)
}

func (this *TeacherController) GetClassList() {
	teacher := this.GetSession("teacher").(models.Teacher)
	teacherService := service.TeacherService{}
	classList, num := teacherService.GetClassList(teacher.Id)
	if num  == 0{
		this.returnJsonWithTeacher(0,"")
	}else{
		this.returnJsonWithTeacher(1,classList)
	}
}

func (this *TeacherController) GetChapterTitle() {
	classId,err := this.GetInt64("classId")
	if err != nil{
		this.returnJson(0,"参数错误")
		return
	}
	teacher := this.GetSession("teacher").(models.Teacher)
	teacherService := service.TeacherService{}
	chapterTitleList,num := teacherService.GetChapterTitle(teacher.Id,classId)
	if num > 0{
		this.returnJsonWithTeacher(1,chapterTitleList)
	}else{
		this.returnJsonWithTeacher(0,1)
	}
}

func (this *TeacherController)Index()  {
	var(
		flag bool
		exec bool
	)
	funArr := this.BaseReqData.GetFun()
	for _, fun := range funArr {
		fun()
		if this.check(){
			return
		}
	}

	exec = false
	this.BaseReqData.Next(func(b bool) {
		flag = b
		exec = true
	},func(_true []func()) {
		if exec {
			if flag{
				for e := range _true {
					_true[e]()
					if this.check(){
						return
					}
				}
				exec = false
			}
		}

	}, func(_false []func()) {
		if exec {
			if !flag {
				for e := range _false {
					_false[e]()
					if this.check(){
						return
					}
				}
				exec = false
			}
		}
	})
	this.defaultRes()
}

//默认返回
func (this *TeacherController)defaultRes()  {
	_,content := this.BaseReqData.ReturnRes()
	this.returnRes(1,content)
}

//检查是否返回
func (this *TeacherController)check()  bool{
	if ok := this.BaseReqData.Check();ok{
		code,content := this.BaseReqData.ReturnRes()
		if ok := this.BaseReqData.IsPlay() ;ok{
			url := this.BaseReqData.GetVideoUrl()
			this.videoPlay(url)
		}else{
			this.returnRes(code,content)
		}
		return true
	}else{
		this.BaseReqData.Clear()
		return false
	}
}


func (this *TeacherController) returnRes(code int8,content interface{}){
	if this.BaseRequestType == "permissions"{
		this.ReturnWithPerSession(code,content)
	}else{
		this.ReturnContent(code,content)
	}
}

//func (this *TeacherController)GetPostData()   {
//	fmt.Println("teacher test")
//}

