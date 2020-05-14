package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)



type AddClass struct {
	TypeId int64
	IsCharge int8
	Name string
	Detail string
	ClassImg string
	Price float64
}
func (this *AddClass) Valid()bool {
	if len(this.Name) == 0{
		return false
	}
	if len(this.Detail) == 0{
		return false
	}
	if len(this.ClassImg) == 0{
		return false
	}
	if this.TypeId == 0 {
		return false
	}
	return true;
}







type ModifyChapterName struct {
	TeacherBaseA
	ClassId int64
	Name string
	teacherId int64
}

//func (this *ModifyChapterName) GetFun()
func (this *ModifyChapterName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{this.Name}
}

func (this *ModifyChapterName) Prepare() {
	teacher := this.ctrl.GetSession("teacher").(models.Teacher)
	this.teacherId = teacher.Id
}


func (this *ModifyChapterName)Model()  {
	r ,_,_ := utils.SqlExec("update chapterTitle set name = ? where id = ? and teacherId = ?",this.Name,this.ClassId,this.teacherId)
	this.setReturn(this.resConver(r),"")
	//*res = utils.SqlTransaction(func(o orm.Ormer,res *bool) {
	//	this.ClassId= 1
	//	utils.SqlTranExec(o,"insert into car(classId) values(?)",res,44)
	//})
}



type DeleteChapterTitle struct {
	TeacherBase
	ChapterTitleId int64
	ClassId int64
}

func (this *DeleteChapterTitle) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ChapterTitleId,this.ClassId},[]string{}
}

func (this *DeleteChapterTitle) GetFun() []func() {
	return []func(){
		func() {
			arr,_ := utils.SqlQuery("select count(*) as count from chapter where chapterTitleId=?",this.ChapterTitleId)
			if arr[0]["count"] != "0"{
				this.setReturn(-1,"该章下还有节，不能删除")
			}
		},
		func() {
			res := utils.SqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"delete  from chapterTitle where id=? and teacherId=? ",res,this.ChapterTitleId,this.Teacher.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `class` set chapterTitleCount = chapterTitleCount - 1 where id = ? and teacherId = ?",res,this.ClassId,this.Teacher.Id)
			})
			this.setReturn(this.resConver(res),"")
		},
	}
}

type AddChapterTitle struct {
	TeacherBase
	ClassId int64
	ChapterTitleName string
	chapterTitleNumStr string
}

func (this *AddChapterTitle) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{this.ChapterTitleName}
}

func (this *AddChapterTitle) GetFun() []func() {

	return []func(){

		func() {
			this.checkBool = utils.SqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into chapterTitle(name,classId,createTime,chapterCount,teacherId) value(?,?,?,?,?)",res,this.ChapterTitleName, this.ClassId,utils.GetTime(),0,this.Teacher.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `class` set chapterTitleCount = chapterTitleCount + 1 where id = ? and teacherId = ?",res,this.ClassId,this.Teacher.Id)
			})
			fmt.Println(this.checkBool)
		},
	}
}

type GetChapterList struct{
	TeacherBaseA
	chapterTitleId string
}

func (this *GetChapterList) GetFun() []func(){
	var(
		chapterList  []models.Chapter
	)
	return []func(){
		func() {
			this.chapterTitleId = this.ctrl.GetString("chapterTitleId")
			this.checkString = this.chapterTitleId
		},
		func() {
			this.sqlQuery("select * from chapter where chapterTitleId=? and teacherId=? ",&chapterList,this.chapterTitleId,this.Teacher.Id)
		},
		func() {
			this.setReturn(1,chapterList)
		},
	}
}

type AddChapter struct {
	TeacherBase
	Time string
	ChapterName string
	ChapterTitleId int64
}

//func (this *AddChapter) validArr() ([]int,[]int64,[]string) {
//	return []int{},[]int64{this.ChapterTitleId},[]string{this.ChapterName,this.Time}
//}
func (this *AddChapter) GetFun() []func() {
	var(
		chapterTitleList []models.ChapterTitle
		lastId int64
		videoUrl string
	)
	return []func(){
		func() {
			this.Time = this.ctrl.GetString("time")
			this.ChapterName = this.ctrl.GetString("ChapterName")
			this.ChapterTitleId,this.checkErr = this.ctrl.GetInt64("chapterTitleId")
			//fmt.Println(this.Time,this.ChapterName,this.ChapterTitleId)
			videoUrl = "static/upload/video/"+utils.GetRandomString(20) +".mp4"
			this.checkErr = this.ctrl.SaveFile(videoUrl,"file")
		},
		func() {
			this.checkZero,this.checkErr = this.sqlQuery("select * from chapterTitle where id=? limit 1",&chapterTitleList,this.ChapterTitleId)
		},
		func() {
			this.judgeNeq(chapterTitleList[0].TeacherId,this.Teacher.Id)
		},
		func() {
			this.checkBool = utils.SqlTransaction(func(o orm.Ormer, res *bool) {
				lastId = utils.SqlTranExec(o,"insert into chapter(classId,name,time,createTime,watchCount,chapterTitleId,teacherId,video)" +
					"values(?,?,?,?,?,?,?,?)",res,chapterTitleList[0].ClassId,this.ChapterName,this.Time,utils.GetTime(),0,this.ChapterTitleId,this.Teacher.Id,videoUrl)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update chapterTitle set chapterCount = chapterCount + 1 where id=?",res,this.ChapterTitleId)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `class` set chapterCount = chapterCount + 1 where id=?",res,chapterTitleList[0].ClassId)
			})
		},
	}
}

type DelChapter struct {
	TeacherBase
	Id int64
	ChapterTitleId int64
	ClassId int64
}

func (this *DelChapter) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id,this.ChapterTitleId,this.ClassId},[]string{}
}

func (this *DelChapter) GetFun() []func() {

	return []func(){
		func() {
			this.checkBool = utils.SqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update chapterTitle set chapterCount = chapterCount - 1 where id=? and teacherId=?",res,this.ChapterTitleId,this.Teacher.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"delete from chapter where id=? and teacherId = ?",res,this.Id,this.Teacher.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update class set chapterCount = chapterCount - 1 where id=? and teacherId=?",res,this.ClassId,this.Teacher.Id)
			})
		},
	}
}

type ModifyChapter struct {
	TeacherBaseA
	Name string
	Time string
	Id int64
}

func (this *ModifyChapter) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{this.Name,this.Time}
}

func (this *ModifyChapter) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("update chapter set name = ?,time = ? where id = ? and teacherId=?",this.Name,this.Time,this.Id,this.Teacher.Id)
		},
	}
}

type GetTypeList struct {
	TeacherBaseA
}

func (this *GetTypeList)GetFun() []func()   {
	var typeList []models.Type
	return []func(){
		func() {
			this.sqlQuery("select id,name from `type`",&typeList)
		},
		func() {
			this.setReturn(1,typeList)
		},
	}
}


type ModifyClass struct {
	TeacherBaseA
	ClassName string
	TypeId int64
	Id int64
	IsCharge int
	IsShelves int
	PicBase64 string
	Detail string
}

func (this *ModifyClass) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.TypeId,this.Id},[]string{this.ClassName,this.Detail}
}

func (this *ModifyClass)GetFun() []func(){
	return []func(){
		func() {
			this.isNotInArr(this.IsCharge,[]interface{}{0,1})
		},
		func() {
			this.isNotInArr(this.IsShelves,[]interface{}{0,1})
		},
	}
}

func (this *ModifyClass)Next(_if func(bool), _true func([]func()),_false func([]func())){
	var picUrl string
	_if(this.PicBase64 != "")
	_true([]func(){
		func() {
			this.checkBool,picUrl = utils.SavePicByBase64(this.PicBase64)
		},
		func() {
			this.update("`class`", map[string]interface{}{
				"name":this.ClassName,
				"isCharge":this.IsCharge,
				"isShelves":this.IsShelves,
				"typeId":this.TypeId,
				"detail":this.Detail,
				"classImg":picUrl,
			}, map[string]interface{}{
				"teacherId":this.Teacher.Id,
				"id":this.Id,
			},"and")
		},
	})

	_false([]func(){
		func() {
			this.update("`class`", map[string]interface{}{
				"name":this.ClassName,
				"isCharge":this.IsCharge,
				"isShelves":this.IsShelves,
				"typeId":this.TypeId,
				"detail":this.Detail,
			}, map[string]interface{}{
				"teacherId":this.Teacher.Id,
				"id":this.Id,
			},"and")
		},
	})
}


type DelClass struct {
	TeacherBaseA
	Id int64
}

func (this *DelClass) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *DelClass)GetFun() []func()   {
	var chapterTitleList []models.ChapterTitle
	return []func(){
		func() {
			this.sqlIsExist("select * from chapterTitle where classId=?",&chapterTitleList,this.Id)
			this.baseStatus = -1
			this.baseContent = "该课程还有章节，不能删除"
		},
		func() {
			this.del("class", map[string]interface{}{
				"id":this.Id,
				"teacherId":this.Teacher.Id,
			},"and")
		},
	}
}

type GetQuestionList struct {
	TeacherBaseA
	ClassId int64
}

func (this *GetQuestionList) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{}
}

func (this *GetQuestionList)GetFun() []func()   {
	type Comm struct {
		models.Communication
		Username string `json:"username" orm:"column(username)"`
		UserImg string `json:"userImg" orm:"column(userImg)"`
		ReplyList map[string]ReplyContent `json:"replyList"`
	}
	var(
		commList []Comm
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		replyContentList map[string]ReplyContent
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where classId=? order by createTime desc",&commList,this.ClassId)
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			this.checkErr  = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			for k, v := range commList {
				v.Username,_ = userMsg[v.UserId]["username"].(string)
				v.UserImg ,_= userMsg[v.UserId]["img"].(string)

				if v.ReplyContent != "" && v.ReplyContent != "{}"{
					this.checkErr = json.Unmarshal([]byte(v.ReplyContent),&replyContentList)
					for key, value := range replyContentList {
						value.PostUsername,_ = userMsg[value.PostUserId]["username"].(string)
						value.SrcUsername,_ = userMsg[value.SrcUserId]["username"].(string)
						replyContentList[key] = value
					}
					v.ReplyList = replyContentList
				}
				commList[k] = v
			}
		},
		func() {
			this.setReturn(1,commList)
		},

	}
}

type TeacherReply struct {
	TeacherBaseA
	CommId int64
	CommIndex string
	Content string
}

func (this *TeacherReply) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.CommId},[]string{this.Content}
}

func (this *TeacherReply)GetFun() []func()   {
	if this.CommIndex == "" {
		return this.replyMaster()
	}else{
		return this.replyOther()
	}

}

func (this *TeacherReply)replyMaster()[]func(){
	var (
		replyContentList map[string]ReplyContent
		commList []models.Communication

		jsons []byte
		replyIndex string
		replyTime int64
		comm models.Communication
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=? limit 1",&commList,this.CommId)
		},
		func() {
			this.checkBool = commList[0].UserId != this.getTeacherId()
			this.baseStatus = -1
			this.baseContent = "不能回复自己"
		},
		func() {
			comm = commList[0]
			if comm.ReplyContentHeader != ""{
				this.checkErr  = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
				replyTime = utils.GetTime()
				reply := ReplyContent{this.getTeacherId(),comm.UserId,this.Content,replyTime,"","","",comm.ReplyContentTail,2,comm.UserType}
				replyIndex = utils.GetRandomString(20)
				replyContentList[replyIndex] = reply

				replyContent := replyContentList[comm.ReplyContentTail]
				replyContent.Next = replyIndex
				replyContentList[comm.ReplyContentTail] = replyContent
				comm.ReplyContentTail = replyIndex
			}else{
				replyContentList = map[string]ReplyContent{}
				replyTime = utils.GetTime()
				reply := ReplyContent{this.getTeacherId(),comm.UserId,this.Content,replyTime,"","","","",2,comm.UserType}
				replyIndex = utils.GetRandomString(20)
				replyContentList[replyIndex] = reply

				comm.ReplyContentTail = replyIndex
				comm.ReplyContentHeader = replyIndex
			}
		},
		func() {
			jsons,this.checkErr = json.Marshal(replyContentList)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,postUserId,postContent,postCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,srcUserType,postUserType) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,comm.UserId,comm.Content,comm.CreateTime,this.getTeacherId(),this.Content,replyTime,0,1,replyIndex,comm.ClassId,comm.Id,comm.UserType,2)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContentTail=?,replyContentHeader=?,replyContent=? where id=?",res,comm.ReplyContentTail,comm.ReplyContentHeader,string(jsons),comm.Id)
			})
		},
	}
}

func (this *TeacherReply)replyOther()[]func() {
	var(
		comm models.Communication
		replyContentList map[string]ReplyContent
		srcReplyContent ReplyContent
		jsonByte []byte
		replyTime int64
		replyIndex string
		commList []models.Communication
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=? limit 1",&commList,this.CommId)
		},
		func() {
			comm = commList[0]
			replyContentStr := comm.ReplyContent
			this.checkErr = json.Unmarshal([]byte(replyContentStr),&replyContentList)
		},
		func() {
			srcReplyContent,this.checkBool = replyContentList[this.CommIndex]
		},
		func() {
			this.checkBool = srcReplyContent.PostUserId != this.getTeacherId()
			this.baseStatus = -1
			this.baseContent = "不能回复自己"
		},
		func() {
			fmt.Println(replyContentList)
			replyTime = utils.GetTime()
			replyIndex = utils.GetRandomString(20)
			tail := replyContentList[comm.ReplyContentTail]
			tail.Next = replyIndex
			replyContentList[comm.ReplyContentTail] = tail

			reply := ReplyContent{this.getTeacherId(),srcReplyContent.PostUserId,this.Content,replyTime,"","","",comm.ReplyContentTail,2,srcReplyContent.PostUserType}
			replyContentList[replyIndex] = reply
			comm.ReplyContentTail = replyIndex
			jsonByte,this.checkErr = json.Marshal(replyContentList)
		},
		func() {
			fmt.Println(srcReplyContent)
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContentTail=?,replyContentHeader=?,replyContent=? where id=?",res,comm.ReplyContentTail,comm.ReplyContentHeader,string(jsonByte),comm.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,postUserId,postContent,postCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,userId,srcUserType,postUserType) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,srcReplyContent.PostUserId,srcReplyContent.Content,srcReplyContent.CreateTime,this.getTeacherId(),this.Content,replyTime,0,1,replyIndex,comm.ClassId,comm.Id,srcReplyContent.SrcUserId,srcReplyContent.PostUserType,2)
			})

		},
	}
}

type SetTeacherName struct {
	TeacherBaseA
	Name string
}

func (this *SetTeacherName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Name}
}

func (this *SetTeacherName)GetFun()[]func() {
	var(
		userMsg map[string]map[string]interface{}
		cacheList []models.Cache
		userMsgByte []byte
	)
	return []func(){
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			userMsg[this.getTeacherId()]["username"] = this.Name
			//userMsg[this.getTeacherId()]["img"] = "/static/img/teacher/default.jpg"
		},
		func() {
			userMsgByte,this.checkErr = json.Marshal(userMsg)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `cache` set userMsg = ? where id = 1",res,string(userMsgByte))
				this.baseContent = "和原来的名字一致"
				this.baseStatus = -1
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update teacher set name =? where id=?",res,this.Name,this.getTeacherId())
			})
		},
		func() {
			this.flushTeacher()
		},
	}
}

type DeleteTeacherComm struct {
	TeacherBaseA
	CommIndex string
	CommId int64
}

func (this *DeleteTeacherComm) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CommId},[]string{}
}

func (this *DeleteTeacherComm)GetFun() []func(){
	if this.CommIndex == ""{
		return this.delMaster()
	}else{
		return this.delOther()
	}
}

func (this *DeleteTeacherComm)delMaster() []func(){
	return []func(){
		func() {

		},
	}
}
func (this *DeleteTeacherComm)delOther() []func(){
	var (
		commList []models.Communication
		replyContentList map[string]ReplyContent
		replyContent ReplyContent
		jsonByte []byte
		comm models.Communication
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=? ",&commList,this.CommId)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(commList[0].ReplyContent),&replyContentList)
		},
		func() {
			replyContent,this.checkBool = replyContentList[this.CommIndex]
		},
		func() {
			this.checkBool = replyContent.PostUserId == this.getTeacherId()
			this.baseContent = "不能删除别人的留言"
			fmt.Println(replyContent.PostUserId,this.getTeacherId())
		},
		func() {
			comm = commList[0]
			if replyContent.Pre == "" && replyContent.Next == ""{
				comm.ReplyContentHeader = ""
				comm.ReplyContentTail = ""
				delete(replyContentList,this.CommIndex)
			}else if replyContent.Pre == "" && replyContent.Next != "" {
				comm.ReplyContentHeader = replyContent.Next
				nextReplyContent := replyContentList[replyContent.Next]
				nextReplyContent.Pre = ""
				replyContentList[replyContent.Next] = nextReplyContent
				delete(replyContentList,this.CommIndex)
			}else if replyContent.Pre != "" && replyContent.Next != ""{
				preReplyContent := replyContentList[replyContent.Pre]
				nextReplyContent := replyContentList[replyContent.Next]
				preReplyContent.Next = replyContent.Next
				nextReplyContent.Pre = replyContent.Pre
				replyContentList[replyContent.Pre] = preReplyContent
				replyContentList[replyContent.Next] = nextReplyContent
				delete(replyContentList,this.CommIndex)
			}else if replyContent.Pre != "" && replyContent.Next == ""{
				comm.ReplyContentTail = replyContent.Pre
				preReplyContent := replyContentList[replyContent.Pre]
				preReplyContent.Next = ""
				replyContentList[replyContent.Pre] = preReplyContent
				delete(replyContentList,this.CommIndex)
			}

		},
		func() {
			delete(replyContentList,this.CommIndex)
			jsonByte ,this.checkErr = json.Marshal(replyContentList)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,postUserId,postCreateTime,srcUserType,postUserType) values(?,?,?,?,?,?,?,?,?,?,?,?)",
					res,replyContent.SrcUserId,replyContent.Content,replyContent.CreateTime,0,2,this.CommIndex,commList[0].ClassId,commList[0].Id,this.getTeacherId(),utils.GetTime(),replyContent.SrcUserType,2)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContent = ? ,replyContentTail=?,replyContentHeader=? where id = ?" ,res,string(jsonByte),comm.ReplyContentTail,comm.ReplyContentHeader,comm.Id)
			})
		},
	}
}

type AddComm struct {
	TeacherBaseA
	Content string
	ClassId int64
}

func (this *AddComm) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{this.Content}
}

func (this *AddComm)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("insert into communication(userId,content,replyContent,classId,createTime,userType) values(?,?,?,?,?,?)",this.getTeacherId(),this.Content,"{}",this.ClassId,utils.GetTime(),2)
		},
	}
}

type TeacherDeleteComm struct {
	TeacherBaseA
	CommId int64
}

func (this *TeacherDeleteComm) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.CommId},[]string{}
}

func (this *TeacherDeleteComm)GetFun() []func(){
	var (
		commList []models.Communication
		comm models.Communication
		replyContentList map[string]ReplyContent
		userIdList []string
		params []interface{}
		notifySlq string
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=?",&commList,this.CommId)
		},
		func() {
			comm =commList[0]
			if comm.UserId != this.getTeacherId() {
				this.checkZero = 0
				this.baseContent = "不可以删除别人的评论"
			}
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
		},
		func() {
			for e := range replyContentList {
				//不向自己发送通知和不重复发通知
				if this.getTeacherId() != replyContentList[e].PostUserId &&!utils.ValIsExistArr(userIdList, replyContentList[e].PostUserId) {
					userIdList = append(userIdList,replyContentList[e].PostUserId)
				}
			}
			//fmt.Println(userIdList)
		},
		func() {
			length := len(userIdList)-1
			values := ""
			time := utils.GetTime()
			for e := range userIdList {
				if e == length{
					values += "(?,?,?,?,?,?,?,?,?)"
				}else{
					values += "(?,?,?,?,?,?,?,?,?),"
				}
				params = append(params,userIdList[e],0,2,comm.ClassId,comm.Id,this.getTeacherId(),time,comm.Content,comm.CreateTime)
			}
			notifySlq = "insert into replyMsgNotify(srcUserId,isHaveRead,notifyType,classId,commId,postUserId,postCreateTime,srcContent,srcCreateTime) values"+values

		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				if len(userIdList) != 0{
					utils.SqlTranExec(o,notifySlq,res,params)
				}
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"delete from communication where id=?",res,this.CommId)
			})
		},
	}
}

type GetNotifyList struct {
	TeacherBaseA
}

func (this *GetNotifyList)GetFun() []func(){
	var (
		replyMsgNotifyList []models.ReplyMsgNotify
		cacheList []models.Cache
		userMsgArr map[string]map[string]interface{}
	)
	return []func(){
		func() {
			this.sqlQuery("select * from replyMsgNotify where srcUserId = ? and isHaveRead=0 order by postCreateTime desc",&replyMsgNotifyList,this.Teacher.Id)
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1 ",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsgArr)
		},
		func() {
			for k, v := range replyMsgNotifyList {
				v.SrcUsername,this.checkBool = userMsgArr[v.SrcUserId]["username"].(string)
				v.SrcUserImg,this.checkBool = userMsgArr[v.SrcUserId]["img"].(string)
				v.PostUsername,this.checkBool = userMsgArr[v.PostUserId]["username"].(string)
				v.PostUserImg,this.checkBool = userMsgArr[v.PostUserId]["img"].(string)
				replyMsgNotifyList[k] = v
			}
			this.check = false
		},
		func() {
			where := "("
			length := len(replyMsgNotifyList)-1
			idList := []int64{}
			for e := range replyMsgNotifyList {
				if length == e{
					where += "?)"
				}else{
					where += "?,"
				}
				idList = append(idList,replyMsgNotifyList[e].Id)
			}
			sql := "update replyMsgNotify set isHaveRead=1 where id in  "+where
			//fmt.Println(sql,idList)
			this.sqlExec(sql,idList)
			this.setReturn(1,replyMsgNotifyList)
		},
	}
}

type GetNotifyDetail struct {
	TeacherBaseA
}

func (this *GetNotifyDetail)GetFun() []func(){
	var(
		commList []models.Communication
		commId int64
		replyContentMap map[string]ReplyContent
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		comm map[string]interface{}
	)
	return []func(){
		func() {
			commId ,this.checkErr= this.ctrl.GetInt64("commId")
			this.checkZero = commId
		},
		func() {
			this.sqlQuery("select * from communication where id = ?",&commList,commId)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(commList[0].ReplyContent),&replyContentMap)
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			comm = map[string]interface{}{
				"username":userMsg[commList[0].UserId]["username"],
				"content":commList[0].Content,
				"createTime":commList[0].CreateTime,
				"classId":commList[0].ClassId,
				"replyContentHeader":commList[0].ReplyContentHeader,
				"replyContentTail":commList[0].ReplyContentTail,
				"userImg":userMsg[commList[0].UserId]["img"],
				"commId":commList[0].Id,
				"userId":commList[0].UserId,
				"userType":commList[0].UserType,
			}
			for k, v := range replyContentMap {
				v.PostUsername,this.checkBool = userMsg[v.PostUserId]["username"].(string)
				v.SrcUsername,this.checkBool = userMsg[v.SrcUserId]["username"].(string)
				replyContentMap[k] = v
			}
			//fmt.Println(userMsg)
			this.check = false
			this.setReturn(1, map[string]interface{}{
				"comm":comm,
				"replyList":replyContentMap,
			})
		},

	}
}


type GetOrderList struct {
	TeacherBaseA
}

func (this *GetOrderList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	type Order struct {
		models.Order
		Username string  `json:"username"`
	}
	var (
		orderList []Order
		start int
		countList []Count
		pageList []utils.Page
		userMsg UserMsg
		cacheList []models.Cache
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
		},
		func() {
			this.sqlQuery("select * from `order` where teacherId =?  order by createTime desc limit ?,10 ",&orderList,this.Teacher.Id,start)
		},
		func() {
			this.sqlQuery("select userMsg from	`cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			for k, v := range orderList {
				v.Username,this.checkBool = userMsg[v.UserId]["username"].(string)
				orderList[k] = v
			}
			this.check = false
		},
		func() {
			this.sqlQuery("select count(*)  as count from `order` where teacherId =?" ,&countList,this.Teacher.Id)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"orderList":orderList,
				"page":pageList,
			})
		},
	}
}


type GetTeacherMessage struct {
	TeacherBaseA
}

func (this *GetTeacherMessage)GetFun() []func(){
	return []func() {
		func() {
			this.flushTeacher()
		},
		func() {
			this.setReturn(1,this.Teacher)
		},
	}
}

type SetTeacherPhone struct {
	TeacherBaseA
	Phone string
}

func (this *SetTeacherPhone) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Phone}
}


func (this *SetTeacherPhone)GetFun() []func(){
	return []func() {
		func() {
			this.sqlExec("update teacher set phone=? where id = ?",this.Phone,this.Teacher.Id);
			this.baseContent = "和原来的电话号码一致"
			this.baseStatus = -1
		},
	}
}

type Withdrawal struct {
	TeacherBaseA
	Username string
	Number string
	Amount float64
	Bank string
}

func (this *Withdrawal) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Username,this.Number,this.Bank}
}

func (this *Withdrawal)GetFun() []func(){
	return []func(){
		func() {
			this.judge(this.Amount <= 0)
		},
		func() {
			this.judge(this.Teacher.Amount < this.Amount)
			this.baseStatus = -1
			this.baseContent = "提现金额不能大于余额"
		},
		func() {
			this.sqlIsExist("select * from withdrawal where teacherId=? and isDo = 0",this.Teacher.Id)
			this.baseStatus = -1
			this.baseContent = "上一次提现未处理，暂不能提现"
		},
		func() {
			var withdrawalId int64
			time := utils.GetTime()
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				withdrawalId = utils.SqlTranExec(o,"insert into withdrawal(username,number,teacherId,createTime,isDo,amount,bank) values(?,?,?,?,?,?,?)",res,
					this.Username,this.Number,this.Teacher.Id,time,0,this.Amount,this.Bank)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update teacher set unMoveAmount=unMoveAmount+?,amount=amount-? where id=?",res,this.Amount,this.Amount,this.Teacher.Id)
			}, func(o orm.Ormer, res *bool) {
				title := "余额提现发起"
				content := "请耐心等待，我们将在两个工作日内处理"
				utils.SqlTranExec(o,"insert into teacherNotify(title,content,teacherId,createTime) values(?,?,?,?)",res,
					title,content,this.Teacher.Id,time)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into platformNotify(title,content,createTime,type,linkId) values(?,?,?,?,?)",
					res,"余额提现发起，请尽快处理","余额提现发起，请尽快处理",time,"withdrawal",withdrawalId)

			})
		},
		func() {
			this.Teacher.UnMoveAmount += this.Amount
			this.Teacher.Amount -= this.Amount
		},
	}
}


type GetTeacherBill struct {
	TeacherBaseA
}

func (this *GetTeacherBill) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		teacherBillList []models.TeacherBill
		start int
		countList []Count
		pageList []utils.Page
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
		},
		func() {
			this.sqlQuery("select * from `teacherBill` where teacherId =?  order by createTime desc limit ?,10 ",&teacherBillList,this.Teacher.Id,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `teacherBill` where teacherId =?" ,&countList,this.Teacher.Id)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"teacherBillList":teacherBillList,
				"page":pageList,
			})
		},
	}
}
type GetTeacherIndexMsg struct {
	TeacherBaseA
}

func (this *GetTeacherIndexMsg)GetFun() []func(){
	var (
		watchRecordList []models.WatchRecord
		yesterdayWatchCount int64
		TeacherNotifyList []models.TeacherNotify
	)
	return []func(){
		func() {
			this.flushTeacher()
		},
		func() {
			yesterday := utils.GetYesterdayData()
			this.sqlQuery("select * from watchRecord where teacherId =? and date = ? limit 1",&watchRecordList,this.Teacher.Id,yesterday)
			this.check = false;
			fmt.Println("yesterday  ",yesterday)
			if len(watchRecordList) == 0 {
				yesterdayWatchCount = 0
			}else{
				yesterdayWatchCount = watchRecordList[0].WatchCount
			}
			this.sqlQuery("select * from teacherNotify where teacherId=? or teacherId=0 order by createTime desc limit 6 ",&TeacherNotifyList,this.Teacher.Id)
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"sellCount":this.Teacher.SellCount,
				"watchCount":this.Teacher.WatchCount,
				"yesterdayWatchCount":yesterdayWatchCount,
				"amount":this.Teacher.Amount,
				"name":this.Teacher.Name,
				"img":this.Teacher.Img,
				"phone":this.Teacher.Phone,
				"notifyList":TeacherNotifyList,
			})
		},
	}
}

type GetTeacherNotifyList struct {
	TeacherBaseA
}
func (this *GetTeacherNotifyList)GetFun() []func(){
	var(
		TeacherNotifyList []models.TeacherNotify
	)
	return []func(){
		func() {
			this.sqlQuery("select * from teacherNotify where teacherId=? or teacherId=0",&TeacherNotifyList,this.Teacher.Id)
		},
		func() {
			fmt.Println(TeacherNotifyList)
		},
	}
}

type TeacherCommPlaints struct {
	TeacherBaseA
	CommId int64
	CommIndex string
}

func (this *TeacherCommPlaints) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CommId},[]string{this.CommIndex}
}

func (this *TeacherCommPlaints)GetFun() []func(){
	var(
		id string
		commComplaintsList []models.CommComplaints
		commList []models.Communication
		comm models.Communication
		replyContentList map[string]ReplyContent
		replayContent ReplyContent
		content string
		postUserId string
		userType int8
		allContent string
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=?",&commList,this.CommId)
		},
		func() {
			comm = commList[0]
			this.checkErr = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
		},
		func() {
			if this.CommIndex == "0"{
				allContent = comm.Content
				postUserId = comm.UserId
				userType = comm.UserType
			}else{
				replayContent,this.checkBool = replyContentList[this.CommIndex]
				allContent  = replayContent.Content
				postUserId = replayContent.PostUserId
				userType = replayContent.PostUserType
			}
			if len(allContent)>60{
				content = utils.SubStrRunes(allContent,20)+"..."
			}else{
				content = allContent
			}

		},
		func() {
			id = strconv.FormatInt(this.CommId,10) + "-" + this.CommIndex
			this.sqlQuery("select id from commComplaints where id =?",&commComplaintsList,id)
			this.check = false
		},
		func() {
			time := utils.GetTime()
			if len(commComplaintsList) == 0{
				this.sqlExec("insert into commComplaints(id,commId,commIndex,content,isDel,createTime,count,userId,classId,userType,allContent) values(?,?,?,?,?,?,?,?,?,?,?)",
					id,this.CommId,this.CommIndex,content,0,time,1,postUserId,comm.ClassId,userType,allContent)
			}else{
				this.sqlExec("update commComplaints set  count = count + 1,createTime = ? where id = ?",time,id)
			}
		},
	}
}


type TeacherGetNotifyList struct {
	TeacherBaseA
}

func (this *TeacherGetNotifyList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		teacherNotifyList []models.TeacherNotify
		start int
		countList []Count
		pageList []utils.Page
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
		},
		func() {
			this.sqlQuery("select * from `teacherNotify` where teacherId= ? or teacherId=0  order by createTime desc limit ?,10 ",&teacherNotifyList,this.Teacher.Id,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `teacherNotify` where teacherId= ? or teacherId=0",&countList,this.Teacher.Id)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"teacherNotifyList":teacherNotifyList,
				"page":pageList,
			})
		},
	}
}

type GetOneTeacherNotify struct {
	TeacherBaseA
	TNId int64
}

func (this *GetOneTeacherNotify) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.TNId},[]string{}
}

func (this *GetOneTeacherNotify) GetFun() []func() {
	var (
		teacherNotifyList []models.TeacherNotify
	)
	return []func(){
		func() {
			this.sqlQuery("select * from teacherNotify where id = ?",&teacherNotifyList,this.TNId)
		},
		func() {
			this.setReturn(1,teacherNotifyList[0])
		},
	}
}

type TeacherGetClassList struct {
	TeacherBaseA
}

func (this *TeacherGetClassList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	type Class struct {
		models.Class
		TypeName string `json:"typeName" orm:"column(typeName)"`
	}
	var (
		classList []Class
		start int
		countList []Count
		pageList []utils.Page
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
		},
		func() {
			this.sqlQuery("select c.*,t.name as typeName from class as c inner join type as t on c.teacherId = ? and t.id=c.typeId order by c.createTime desc limit ?,10" ,&classList,this.Teacher.Id,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `class` where teacherId =?" ,&countList,this.Teacher.Id)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"classList":classList,
				"page":pageList,
			})
		},
	}
}





//type TeacherGetClassList struct {
//	TeacherBaseA
//}

//func (this *TeacherGetClassList) GetFun() []func() {
//	type Class struct {
//		models.Class
//		TypeName string `json:"typeName" orm:"column(typeName)"`
//	}
//	var(
//		classList []Class
//	)
//	return []func(){
//		func() {
//			this.sqlQuery("select c.*,t.name as typeName from class as c inner join type as t on c.teacherId = ? and t.id=c.typeId " ,&classList,this.Teacher.Id)
//		},
//		func() {
//			//this.setReturn(1, map[string]interface{})
//		},
//	}
//}

type TeacherGetTypeList struct {
	TeacherBaseA
}
func (this *TeacherGetTypeList) GetFun() []func() {
	var (
		typeList []models.Type
	)
	return []func(){
		func() {
			this.sqlQuery("select * from type ",&typeList)
		},
		func() {
			this.setReturn(1,typeList)
		},
	}
}

type TeacherLogin struct {
	TeacherBaseA
	Phone string
	Password string

}
func (this *TeacherLogin) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.Phone,this.Password}
}
func (this *TeacherLogin) GetFun() []func() {
	var(
		teacherList []models.Teacher
	)
	return []func(){
		func() {
			this.Password = utils.Md5(this.Password)
			this.sqlQuery("select * from teacher where phone = ? and password =?",&teacherList,this.Phone,this.Password)
		},
		func() {
			this.checkBool = teacherList[0].IsBan == 0
			this.errMsg("你已经被管理员禁止登录")
		},
		func() {
			this.ctrl.SetSession("teacher",teacherList[0])
		},
	}
}

type RegisterCode struct {
	Code string
	Time int64
	Phone string
}

type TeacherGetCode struct {
	TeacherBaseA
	Phone string
}
func (this *TeacherGetCode) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.Phone}
}

func (this *TeacherGetCode) GetFun() []func() {
	var code string
	return []func(){
		func() {
			this.sqlIsExist("select * from teacher where phone=?", this.Phone)
			this.errMsg(this.Phone + "号码已注册")
		},
		func() {
			code = strconv.FormatInt(utils.RandNum(1000,9999),10)
			//fmt.Println(code)
			this.checkErr = utils.SendSms("LTAI4Fiop5swEDrRhUD5vo1B", "xImuBjvIQjWvcrvFDLM5IdSD5if51A", this.Phone, "萤火虫", "{\"code\":"+code+"}", "SMS_124460061");
		},
		func() {
			time := utils.GetTime() + 60 * 30
			registerCode := RegisterCode{code,time,this.Phone}
			this.ctrl.SetSession("code",registerCode)
		},
	}
}

type TeacherGetCodeResetPwd struct {
	TeacherBaseA
	Phone string
}
func (this *TeacherGetCodeResetPwd) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.Phone}
}

func (this *TeacherGetCodeResetPwd) GetFun() []func() {
	var(
		code string
		teacherList []models.Teacher
	)
	return []func(){
		func() {
			this.sqlQuery("select * from teacher where phone=? limit 1", &teacherList,this.Phone)
			this.errMsg(this.Phone + "号码未注册")
		},
		func() {
			code = strconv.FormatInt(utils.RandNum(1000,9999),10)
			fmt.Println(code)
			//this.checkErr = utils.SendSms("LTAI4Fiop5swEDrRhUD5vo1B", "xImuBjvIQjWvcrvFDLM5IdSD5if51A", this.Phone, "萤火虫", "{\"code\":"+code+"}", "SMS_124380039");
		},
		func() {
			time := utils.GetTime() + 60 * 30
			registerCode := RegisterCode{code,time,this.Phone}
			this.ctrl.SetSession("codeResetPwd",registerCode)
		},
	}
}

type TeacherRegister struct {
	TeacherBaseA
	Username string
	Password string
	Code string
}

func (this *TeacherRegister) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.Username,this.Password,this.Code}
}

func (this *TeacherRegister) GetFun() []func() {
	var(
		teacherList []models.Teacher
		code RegisterCode
		time int64
		cacheList []models.Cache
		userMsg UserMsg
		teacherId int64
		img string
		userMsgStr string
	)
	return []func(){
		func() {
			code = this.ctrl.GetSession("code").(RegisterCode)
			time = utils.GetTime()
			this.checkBool = code.Time > time
			this.errMsg("验证码已经过期")
		},
		func() {
			this.checkBool = code.Code == this.Code
			this.errMsg("验证码不正确")
		},
		func() {
			this.sqlIsExist("select * from teacher where phone=?",code.Phone)
			this.errMsg(code.Phone + "号码已注册")
		},
		func() {
			this.sqlQuery("select * from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			img = "/static/img/teacher/default.jpg"
			this.Password = utils.Md5(this.Password)
			teacherId = this.sqlExec("insert into teacher(name,phone,password,amount,createTime,img,unMoveAmount,sellCount,watchCount,isBan) values(?,?,?,?,?,?,?,?,?,?)",
				this.Username,code.Phone,this.Password,0,time,img,0,0,0,0)
		},
		func() {
			userMsg[strconv.FormatInt(teacherId,10)] = map[string]interface{}{
				"username":this.Username,
				"img":img,
			}
		},
		func() {
			var b []byte
			b,this.checkErr = json.Marshal(userMsg)
			userMsgStr = string(b)
		},
		func() {
			this.sqlExec("update `cache` set userMsg = ? where id=1",userMsgStr)
		},
		func() {
			this.ctrl.SetSession("code",nil)
		},
		func() {
			this.sqlQuery("select * from teacher where phone = ? limit 1",&teacherList,code.Phone)
		},
		func() {
			this.ctrl.SetSession("teacher",teacherList[0])
		},
	}
}

type TeacherLogout struct {
	TeacherBaseA
}
func (this *TeacherLogout) GetFun() []func() {
	return []func(){
		func() {
			this.ctrl.SetSession("teacher",nil)
		},
	}
}

type TeacherResetPwd struct {
	TeacherBaseA
	Password string
	Code string
}

func (this *TeacherResetPwd) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Password,this.Code}
}
func (this *TeacherResetPwd)GetFun() []func(){
	var(
		code RegisterCode
		time int64
		teacherlist []models.Teacher
	)
	return []func(){
		func() {
			this.checkBool = this.ctrl.GetSession("codeResetPwd") != nil

		},
		func() {
			code = this.ctrl.GetSession("codeResetPwd").(RegisterCode)
		},
		func() {
			time = utils.GetTime()
			this.checkBool = code.Time >  time
			this.errMsg("验证码已过期")
		},
		func() {
			this.checkBool = code.Code == this.Code
			this.errMsg("验证码不正确")
		},
		func() {
			this.Password = utils.Md5(this.Password)
			this.sqlExec("update teacher set password = ? where phone = ? limit 1",this.Password, code.Phone)
		},
		func() {
			this.sqlQuery("select * from teacher where phone =?",&teacherlist,code.Phone)
		},
		func() {
			this.ctrl.SetSession("codeResetPwd",nil)
			this.ctrl.SetSession("teacher",teacherlist[0])
		},
	}
}

type UpdateIsShelves struct {
	TeacherBaseA
	Status int
	ClassId int64
}

func (this *UpdateIsShelves) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{}
}
func (this *UpdateIsShelves)GetFun() []func(){
	var(
		classList []models.Class
	)
	return []func(){
		func() {
			this.sqlQuery("select * from class where id=? limit 1",&classList,this.ClassId)
		},
		func() {
			if this.Status == 1{
				this.checkBool =  classList[0].IsBan == 0
				this.errMsg("该课程已经被管理员禁止上架")
			}

		},
		func() {
			this.sqlExec("update class set isShelves = ? where id = ? and teacherId = ?",this.Status,this.ClassId,this.Teacher.Id)
		},
	}
}