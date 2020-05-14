package models

import (
	"cloudClass/utils"
	"github.com/astaxie/beego/orm"
)

type TeacherModel struct {

}

func (this *TeacherModel) GetTeacherById(id int64) ([]orm.Params,int64){
	//var maps []orm.Params
	//o :=orm.NewOrm()
	//num ,_ := o.Raw("select * from teacher where id =? limit 1",8).Values(&maps)
	//return maps, num
	return utils.SqlQuery("select * from teacher where id =? limit 1", id)
}

func (this * TeacherModel)GetTeacherByPhone(phone string)([]orm.Params,int)  {
	var maps []orm.Params
	o := orm.NewOrm()
	num ,_ := o.Raw("select * from teacher where phone=? limit 1", phone).Values(&maps)
	return maps,int(num)
}

func (this *TeacherModel) Add(phone string, password string,amount int,createTime int64) bool{
	o := orm.NewOrm()
	_,err := o.Raw("insert into teacher(phone,password,amount,createTime) values(?,?,?,?)",phone,password,amount,createTime).Exec()
	if err == nil{
		return true
	}
	return false
}

//func (this *TeacherModel) exec(sql string,params ... interface{})(isSuccess bool,lastId int64,affected int64) {
//	o := orm.NewOrm()
//	res,err := o.Raw(sql,params).Exec()
//	if err != nil{
//		isSuccess = false
//		return
//	}
//	isSuccess = true
//	lastId,_ = res.LastInsertId()
//	affected,_ = res.RowsAffected()
//	return
//}

func (this *TeacherModel) AddClass(class Class) (isSuccess bool,lastId int64,affected int64){
	return utils.SqlExec("insert into `class`(name,isCharge,createTime,teacherId,typeId,classImg,watchCount,detail,isShelves,chapterCount,chapterTitleCount,price,star) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		class.Name,class.IsCharge,class.CreateTime,class.TeacherId,class.TypeId,class.ClassImg,class.WatchCount,class.Detail,class.IsShelves,class.ChapterCount,class.ChapterTitleCount,class.Price,class.Star)
}

/**
	返回值 第一个是返回查询到的数据列表，第二个是一共有几条数据，
	判断是否查询到数据，只要判断第二个值是否大于0即可
 */
func (this *TeacherModel) GetClassType()([]orm.Params,int64) {
	return utils.SqlQuery("select id,name from `type`")
}

func (this *TeacherModel) GetClassList(teacherId int64) ([]orm.Params,int64){
	return utils.SqlQuery("select c.*,t.name as typeName from `class` as c  inner join `type` as t on c.teacherId = ? and t.id = c.typeId",teacherId)
}

func (this *TeacherModel) GetChapterTitleByClassId(id interface{}) ([]orm.Params,int64){
	return utils.SqlQuery("select id,chapterCount,count(*)  from `chapterTitle` where classId = ? ",id)
}

func (this *TeacherModel) GetChapterTitleByCIdAndTId(teacherId int64, classId int64) ([]orm.Params,int64){
	return utils.SqlQuery("select * from chapterTitle where classId = ? and teacherId = ? ",classId,teacherId)
}

func (this *TeacherModel) GetChapterTitleCount(classId int64, teacherId int64)([]orm.Params,int64) {
	return utils.SqlQuery("select chapterTitle from `class` where id=? and teacherId = ?",classId,teacherId)
}

