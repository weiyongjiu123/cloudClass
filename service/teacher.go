package service

import (
	"cloudClass/models"
	"cloudClass/request"
	"cloudClass/utils"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type TeacherService struct {
	teacherModel models.TeacherModel
}

func (this *TeacherService)ServiceInit()  {
	this.teacherModel = models.TeacherModel{}
}

func (this *TeacherService)Login (phone string ,password string) (isLogin bool,teacher orm.Params){
	teacherModel := models.TeacherModel{}
	maps, num := teacherModel.GetTeacherByPhone(phone)
	if(num < 1) {
		isLogin = false
		return
	}
	teacher = maps[0]
	if(utils.Md5(password) == teacher["password"]){
		isLogin =  true
	}else{
		isLogin = false
	}
	return
}

func (this *TeacherService) IsHaveUser(phone string)  bool{
	teacherModel := models.TeacherModel{}
	_, num := teacherModel.GetTeacherByPhone(phone)
	if num > 0{
		return true
	}
	return false
}

func (this *TeacherService) Register(phone string, password string) (bool,map[string]interface{}){
	teacherModel := models.TeacherModel{}
	teacherMap := make(map[string]interface{})
	time := utils.GetTime()
	password = utils.Md5(password)
	teacherMap["phone"] = phone
	teacherMap["password"] = password
	teacherMap["amount"] = 0
	teacherMap["createTime"] = time
	return teacherModel.Add(phone,password,0,time),teacherMap
}

func (this *TeacherService) AddClass(data request.AddClass,teacherId int64) bool{
	teacherModel := models.TeacherModel{}
	//isCharge,_ :=strconv.Atoi(data["isCharge"])
	//typeId,_ :=strconv.Atoi(data["typeId"])
	tableClass := models.Class{
		0,
		data.Name,
		data.Detail,
		data.IsCharge,
		utils.GetTime(),
		teacherId,
		data.TypeId,
		data.ClassImg,
		0,
		0,
		0,
		0,
		data.Price,
		0,
		4.5,
		0,
		"{}",
	}
	addRes,_,_ :=teacherModel.AddClass(tableClass)
	return addRes
}

func (this *TeacherService) GetTeacherById(id int64)([]orm.Params,int64) {
	teacherModel := models.TeacherModel{}
	return  teacherModel.GetTeacherById(id)
}

func (this *TeacherService) GetClassType() ([]orm.Params,int64){
	teacherModel := models.TeacherModel{}
	return teacherModel.GetClassType()
}

func (this *TeacherService) GetClassList(teacherId int64) ([]orm.Params,int64){
	teacherModel := models.TeacherModel{}
	classList,num :=  teacherModel.GetClassList(teacherId)
	for _, v := range classList {
		time,_ := strconv.ParseInt(v["createTime"].(string), 10,64)
		v["createTime"] = utils.GetDate(time)
	}
	return classList,num
}

func (this *TeacherService) GetChapterTitle(teacherId int64, classId int64) ([]orm.Params,int64){
	teacherModel := models.TeacherModel{}
	return teacherModel.GetChapterTitleByCIdAndTId(teacherId,classId)
}

func (this *TeacherService) GetChapterTitleCount(classId int64, teacherId int64) {
	teacherModel := models.TeacherModel{}
	teacherModel.GetChapterTitleCount(classId,teacherId)
}

