package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type DoWithdrawal struct {
	AdminA
	Id int64
}

func (this *DoWithdrawal) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *DoWithdrawal) GetFun() []func() {
	return []func() {
		func() {
			this.sqlExec("update withdrawal set isDo = 1 where id = ?",this.Id)
		},
	}
}

type AdminGetIndexMsg struct {
	AdminA
}
func (this *AdminGetIndexMsg) GetFun() []func() {
	type Teacher struct {
		SellCount int64 `json:"sellCount" orm:"column(sellCount)"`
		WatchCount int64  `json:"watchCount" orm:"column(watchCount)"`
		Count int64  `json:"count" orm:"column(count)"`
	}
	type WatchRecord struct {
		WatchCount int64 `json:"watchCount" orm:"column(watchCount)"`
	}
	type User struct {
		Count int64   `json:"count" orm:"column(count)"`
	}
	var (
		teacherList []Teacher
		watchRecordList []WatchRecord
		userList []User
		platformNotifyList []models.PlatformNotify
	)
	return []func(){
		func() {
			this.sqlQuery("select sum(sellCount) as sellCount,sum(watchCount) as watchCount,count(*) as count from teacher ",&teacherList)
			this.sqlQuery("select sum(watchCount) as watchCount  from watchRecord where date = ? ",&watchRecordList,utils.GetYesterdayData())
			this.sqlQuery("select count(*) as count from `user`",&userList)
			this.sqlQuery("select * from platformNotify order by createTime desc  limit 10",&platformNotifyList)
			this.check = false
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"sellCount":teacherList[0].SellCount,
				"watchCount":teacherList[0].WatchCount,
				"yesterdayWatchCount":watchRecordList[0].WatchCount,
				"userCount":userList[0].Count,
				"teacherCount":teacherList[0].Count,
				"platformNotifyList":platformNotifyList,
			})
		},
	}
}

type AdminGetUserList struct {
	AdminA
}



func (this *AdminGetUserList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		userList []models.User
		start int
		countList []Count
		pageList []utils.Page
		sortId int
		sortList = []string{"createTime","amount"}
		keyword string
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			sortId = this.ctrl.GetInt("sortId")
			keyword = this.ctrl.GetString("keyword")
		},
		func() {
			if keyword == ""{
				this.sqlQuery("select * from `user` order by "+sortList[sortId]+" desc limit ?,10 ",&userList,start)
			}else{
				this.sqlQuery("select * from `user` where id = ?",&userList,keyword)
			}
		},
		func() {
			if keyword == ""{
				this.sqlQuery("select count(*)  as count from `user`",&countList)
			}else{
				this.sqlQuery("select count(*)  as count from `user` where id = ?",&countList,keyword)
			}

			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"userList":userList,
				"page":pageList,
			})
		},
	}
}

type UserBan struct {
	AdminA
	IsBan int
	UserId string
}
func (this *UserBan) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.UserId}
}
func (this *UserBan) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("update `user` set isBan=? where id = ?",this.IsBan,this.UserId)
		},
	}
}

type GetTeacherList struct {
	AdminA
}


func (this *GetTeacherList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		teacherLisst []models.Teacher
		start int
		countList []Count
		pageList []utils.Page
		sortId int
		sortList = []string{"createTime","amount","unMoveAmount","sellCount","watchCount"}
		keyword string
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			sortId = this.ctrl.GetInt("sortId")
			keyword = this.ctrl.GetString("keyword")
		},
		func() {
			if keyword == "" {
				this.sqlQuery("select * from `teacher` order by "+sortList[sortId]+" desc limit ?,10",&teacherLisst,start)
			}else{
				this.sqlQuery("select * from teacher where id= ?",&teacherLisst,keyword)
			}
		},
		func() {
			if keyword == "" {
				this.sqlQuery("select count(*)  as count from `teacher`",&countList)
			}else{
				this.sqlQuery("select count(*)  as count from `teacher` where id = ?",&countList,keyword)
			}
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"teacherList":teacherLisst,
				"page":pageList,
			})
		},
	}
}

type TeacherBan struct {
	AdminA
	IsBan int
	TeacherId int64
}
func (this *TeacherBan) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.TeacherId},[]string{}
}
func (this *TeacherBan) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("update `teacher` set isBan=? where id = ?",this.IsBan,this.TeacherId)
		},
	}
}

type GetComplaintsList struct {
	AdminA
}

func (this *GetComplaintsList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		complaintsList []models.Complaints
		start int
		countList []Count
		pageList []utils.Page
		isDo int
		isDoMap = []string{" where isDo = 0"," where isDo = 1",""}
		paramsId int
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			isDo = this.ctrl.GetInt("isDo")
			paramsId = this.ctrl.GetInt("id")
		},
		func() {
			if paramsId > 0 {
				this.sqlQuery("select * from `complaints` where id=? ",&complaintsList,paramsId)
				pageList = utils.GetPage(1,start)
				this.setReturn(1,map[string]interface{}{
					"complaintsList":complaintsList,
					"page":pageList,
				})
			}
		},
		func() {
			this.sqlQuery("select * from `complaints` "+isDoMap[isDo]+" order by createTime desc limit ?,10 ",&complaintsList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `complaints`" +isDoMap[isDo],&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"complaintsList":complaintsList,
				"page":pageList,
			})
		},
	}
}

type ComplaintsDo struct {
	AdminA
	IsDo int8
	Id int64
}

func (this *ComplaintsDo) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *ComplaintsDo) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("update complaints set isDo = ? where id = ?",this.IsDo,this.Id)
		},
	}
}



type GetWithdrawalList struct {
	AdminA
}

func (this *GetWithdrawalList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		withdrawalList []models.Withdrawal
		start int
		countList []Count
		pageList []utils.Page
		isDo int
		isDoMap = []string{" where isDo = 0"," where isDo > 0",""}
		id int
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			isDo = this.ctrl.GetInt("isDo")
			id = this.ctrl.GetInt("id")
		},
		func() {
			if id > 0{
				this.sqlQuery("select * from `withdrawal` where id =?",&withdrawalList,id)
				pageList = utils.GetPage(1,start)
				this.setReturn(1,map[string]interface{}{
					"withdrawalList":withdrawalList,
					"page":pageList,
				})
			}
		},
		func() {
			this.sqlQuery("select * from `withdrawal` "+isDoMap[isDo]+" order by createTime desc limit ?,10 ",&withdrawalList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `withdrawal`" +isDoMap[isDo],&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"withdrawalList":withdrawalList,
				"page":pageList,
			})
		},
	}
}

type WithdrawalDo struct {
	AdminA
	IsDo int
	Id int64
	FailReason string
}

func (this *WithdrawalDo) validArr() ([]int,[]int64,[]string) {
	return []int{this.IsDo},[]int64{this.Id},[]string{}
}

func (this *WithdrawalDo) GetFun() []func() {
	var (
		withdrawalList []models.Withdrawal
		notifyContent string
		notifyTitle string
		withdrawal models.Withdrawal
	)
	return []func(){
		func() {
			this.sqlQuery("select * from withdrawal where id = ?",&withdrawalList,this.Id)
		},
		func() {
			withdrawal = withdrawalList[0]
			if this.IsDo == 2{
				this.checkString = this.FailReason
				notifyContent = this.FailReason
				notifyTitle = "提现失败"
			}else{

				notifyContent = "您的"+ utils.FloatToString(withdrawal.Amount) + "￥提现成功"
				notifyTitle = "提现成功"
			}
		},
		func() {
			time := utils.GetTime()
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update withdrawal set isDo = ?,doTime = ? where id=?",res,this.IsDo,time,this.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into teacherNotify(title,content,teacherId,createTime) values(?,?,?,?)",res,notifyTitle,notifyContent,withdrawal.TeacherId,time)
			}, func(o orm.Ormer, res *bool) {
				if this.IsDo == 2{
					utils.SqlTranExec(o,"update teacher set amount = amount + ?,unMoveAmount = unMoveAmount - ? where id=?",res,withdrawal.Amount,withdrawal.Amount,withdrawal.TeacherId)
				}else{
					utils.SqlTranExec(o,"update teacher set unMoveAmount = unMoveAmount - ? where id = ?",res,withdrawal.Amount,withdrawal.TeacherId)
				}
			}, func(o orm.Ormer, res *bool) {
				id  := strconv.FormatInt(utils.GetTime(),10) + strconv.FormatInt(utils.RandNum(1000000000,9999999999),10)
				description := utils.FloatToString(withdrawal.Amount) + "￥提现成功"
				if this.IsDo == 1{
					utils.SqlTranExec(o,"insert into teacherBill(description,amount,createTime,orderId,teacherId) values(?,?,?,?,?)",res,
						description,-withdrawal.Amount,time,id,withdrawal.TeacherId)
				}
			})
		},
	}
}

type GetClassList struct {
	AdminA
}


func (this *GetClassList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	type Class struct {
		models.Class
		TeacherName string `json:"teacherName"`
		TypeName string `json:"typeName"`
	}
	type ClassType struct {
		Id int64 `json:"id" orm:"column(id)"`
		Name string `json:"name" orm:"column(name)"`
	}
	var (
		classList []Class
		start int
		countList []Count
		pageList []utils.Page
		sortId int
		sortList = []string{"createTime","watchCount","chapterTitleCount","chapterCount","price","star"}
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		classTypeList []ClassType
		keyword string
		classId string
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			sortId = this.ctrl.GetInt("sortId")
			keyword = this.ctrl.GetString("keyword")
			classId = this.ctrl.GetString("classId")
		},
		func() {
			this.sqlQuery("select * from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			this.sqlQuery("select * from type",&classTypeList)
		},
		func() {
			if classId != ""{
				this.sqlQuery("select * from class where id = ?",&classList,classId)
				if len(classList) == 0{
					this.setReturn(0,"")
				}else{
					classList[0].TeacherName = userMsg[strconv.FormatInt(classList[0].TeacherId,10)]["username"].(string)
					for _, v := range classTypeList {
						if v.Id == classList[0].TypeId{
							classList[0].TypeName = v.Name
						}
					}
					this.setReturn(1,map[string]interface{}{
						"classList":classList,
						"page":pageList,
					})
				}

			}
		},
		func() {
			if keyword == ""{
				this.sqlQuery("select * from `class` order by "+sortList[sortId]+" desc limit ?,10",&classList,start)
			}else{
				this.sqlQuery("select * from `class` where name like '%"+keyword+"%'",&classList)
			}
		},

		func() {
			for k, v := range classList {
				v.TeacherName,this.checkBool = userMsg[strconv.FormatInt(v.TeacherId,10)]["username"].(string)
				for _, value := range classTypeList {
					if value.Id == v.TypeId{
						v.TypeName = value.Name
					}
				}
				classList[k] = v
			}
			if keyword == ""{
				this.sqlQuery("select count(*)  as count from `class`",&countList)
			}else{
				this.sqlQuery("select count(*)  as count from `class` where name like '%"+keyword+"%'",&countList)
			}
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

type ClassBan struct {
	AdminA
	IsBan int
	ClassId int64
	BanReason string
}

func (this *ClassBan) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{}
}

func (this *ClassBan) GetFun() []func() {
	var (
		classList []models.Class
		notifyContent string
		notifyTitle string
		class models.Class
	)
	return []func(){
		func() {
			this.sqlQuery("select * from class where id = ?",&classList,this.ClassId)
		},
		func() {
			class = classList[0]
			if this.IsBan == 1{
				this.checkString = this.BanReason
				notifyContent = "禁止上架原因："+this.BanReason
				notifyTitle = "《"+class.Name+"》" + "课堂禁止上架"
			}else{
				notifyContent = "《"+class.Name+"》" + "课堂已解除禁止上架，请手动上架该课程"
				notifyTitle =  "《"+class.Name+"》" + "课堂解除解除上架"
			}
		},
		func() {
			time := utils.GetTime()
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				if this.IsBan == 1{
					utils.SqlTranExec(o,"update class set isBan = 1,isShelves = 0 where id=?",res,this.ClassId)
				}else{
					utils.SqlTranExec(o,"update class set isBan = 0 where id=?",res,this.ClassId)
				}
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into teacherNotify(title,content,teacherId,createTime) values(?,?,?,?)",res,notifyTitle,notifyContent,class.TeacherId,time)
			})
		},
	}
}


type AdminGetOrderList struct {
	AdminA
}


func (this *AdminGetOrderList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		orderList []models.Order
		start int
		countList []Count
		pageList []utils.Page
		sortId int
		sortList = []string{"createTime","amount"}
		keyword string
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			sortId = this.ctrl.GetInt("sortId")
			keyword = this.ctrl.GetString("keyword")
		},
		func() {
			if keyword == "" {
				this.sqlQuery("select * from `order` order by "+sortList[sortId]+" desc limit ?,10",&orderList,start)
			}else{
				this.sqlQuery("select * from `order` where orderId= ?",&orderList,keyword)
			}
		},
		func() {
			if keyword == "" {
				this.sqlQuery("select count(*)  as count from `order`",&countList)
			}else{
				this.sqlQuery("select count(*)  as count from `order` where orderId = ?",&countList,keyword)
			}
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

type GetBillOne struct {
	AdminA
}

func (this *GetBillOne) GetFun() []func() {
	var(
		billList []models.Bill
		id string
	)
	return []func(){
		func() {
			id = this.ctrl.GetString("id")
			this.checkString = id
		},
		func() {
			this.sqlQuery("select * from bill where id = ? limit 1",&billList,id)
		},
		func() {
			this.setReturn(1,billList[0])
		},
	}
}


type AdminGetQuestionList struct {
	AdminA
}


func (this *AdminGetQuestionList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		questionList []models.Question
		start int
		countList []Count
		pageList []utils.Page
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
		},
		func() {
			this.sqlQuery("select * from `question` order by createTime desc  limit ?,10",&questionList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `question`",&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"questionList":questionList,
				"page":pageList,
			})
		},
	}
}

type AddQuestion struct {
	AdminA
	Question string
	Answer string
}

func (this *AddQuestion) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Question,this.Answer}
}
func (this *AddQuestion) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("insert into question(question,answer,createTime)  values(?,?,?)",this.Question,this.Answer,utils.GetTime())
		},
	}
}

type QuestionDel struct {
	AdminA
	Id int64
}
func (this *QuestionDel) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *QuestionDel) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("delete from question where id = ?",this.Id)
		},
	}
}

type AdminGetTypeList struct {
	AdminA
}

func (this *AdminGetTypeList) GetFun() []func() {
	var(
		typeList []models.Type

	)
	return []func(){
		func() {
			this.sqlQuery("select * from `type` ",&typeList)
		},
		func() {
			this.setReturn(1,typeList)
		},
	}
}

type ModifyTypeName struct {
	AdminA
	Name string
	Id int64
}
func (this *ModifyTypeName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{this.Name}
}

func (this *ModifyTypeName) GetFun() []func() {
	return []func(){
		func() {
			this.sqlIsExist("select * from `type` where name=?",this.Name)
			this.baseStatus = -1;
			this.baseContent = "已有相同类型名字"
		},
		func() {
			this.sqlExec("update `type` set name=? where id = ?",this.Name,this.Id)
		},
	}
}

type AddType struct {
	AdminA
	Name string
}

func (this *AddType) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Name}
}

func (this *AddType) GetFun() []func() {
	return []func(){
		func() {
			this.sqlIsExist("select * from `type` where name=?",this.Name)
			this.baseStatus = -1;
			this.baseContent = "已有相同类型名字"
		},
		func() {
			this.sqlExec("insert into `type`(name,createTime) value(?,?)",this.Name,utils.GetTime())
		},
	}
}

type DelType struct {
	AdminA
	Id int64
}

func (this *DelType) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *DelType) GetFun() []func() {
	return []func(){
		func() {
			this.sqlIsExist("select * from `class` where typeId=? limit 1",this.Id)
			this.baseStatus = -1
			this.baseContent = "该类型下还有课程，不能删除"
		},
		func() {
			this.sqlExec("delete from `type` where id = ?",this.Id)
		},
	}
}

type GetAdList struct {
	AdminA
}

func (this *GetAdList) GetFun() []func() {
	var(
		adList []models.Ad
	)
	return []func(){
		func() {
			this.sqlQuery("select * from ad",&adList)
		},
		func() {
			this.setReturn(1,adList)
		},
	}
}

type AddAd struct {
	AdminA
}
func (this *AddAd) GetFun() []func() {
	var(
		picUrl string
		classId string
		classList []models.Class
	)
	return []func(){
		func() {
			classId = this.ctrl.GetString("classId")
			this.checkString = classId
		},
		func() {
			this.sqlQuery("select * from class where id=?",&classList,classId)
			this.baseStatus = -1;
			this.baseContent = "课程id'"+classId+"’不存在"
		},
		func() {
			picUrl = "static/upload/ad/"+utils.GetRandomString(20) +".jpg"
			this.checkErr = this.ctrl.SaveFile(picUrl,"file")
		},
		func() {
			this.sqlExec("insert into ad(picUrl,classId,createTime) values(?,?,?)","/"+picUrl,classId,utils.GetTime())
		},
	}
}

type DelAd struct {
	AdminA
	Id int64
}

func (this *DelAd) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *DelAd) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("delete from ad where id=?",this.Id)
		},
	}
}

type AdminGetQuestion struct {
	AdminA
	ClassId int64
}

func (this *AdminGetQuestion) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{}
}

func (this *AdminGetQuestion)GetFun() []func()   {
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


type GetCommPlaintList struct {
	AdminA
}

func (this *GetCommPlaintList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		comCommplaintList []models.CommComplaints
		start int
		countList []Count
		pageList []utils.Page
		isDo int
		isDoMap = []string{" where isDel = 0"," where isDel = 1",""}
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			isDo = this.ctrl.GetInt("isDo")
		},
		func() {
			this.sqlQuery("select * from `commComplaints` "+isDoMap[isDo]+" order by createTime desc limit ?,10 ",&comCommplaintList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `commComplaints`" +isDoMap[isDo],&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"commComplaintsList":comCommplaintList,
				"page":pageList,
			})
		},
	}
}

type AdminDeleteComm struct {
	AdminA
	CommId int64
	ComplaintsId string
}

func (this *AdminDeleteComm) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.ComplaintsId}
}

func (this *AdminDeleteComm)GetFun() []func(){
	var (
		commList []models.Communication
		comm models.Communication
		replyContentList map[string]ReplyContent
		userIdList []string
		params []interface{}
		notifySlq string
		complaintsList []models.CommComplaints
	)
	return []func(){
		func() {
			this.sqlQuery("select * from commComplaints where id=?",&complaintsList,this.ComplaintsId)
		},
		func() {
			this.checkBool = complaintsList[0].IsDel == 0
			this.CommId = complaintsList[0].CommId
		},
		func() {
			this.sqlQuery("select * from communication where id=?",&commList,this.CommId)
			this.errMsg("该留言已被用户自己删除")
		},
		func() {
			comm =commList[0]
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
		},
		func() {
			userIdList = append(userIdList,comm.UserId)
			for e := range replyContentList {
				//不向自己发送通知和不重复发通知
				if !utils.ValIsExistArr(userIdList, replyContentList[e].PostUserId) {
					userIdList = append(userIdList,replyContentList[e].PostUserId)
				}
			}
		},
		func() {
			length := len(userIdList)-1
			values := ""
			time := utils.GetTime()
			postContent := "管理员删除了这条信息"
			for e := range userIdList {
				if e == length{
					values += "(?,?,?,?,?,?,?,?,?,?,?)"
				}else{
					values += "(?,?,?,?,?,?,?,?,?,?,?),"
				}
				params = append(params,userIdList[e],0,3,comm.ClassId,comm.Id,0,time,comm.Content,comm.CreateTime,postContent,3)
			}
			notifySlq = "insert into replyMsgNotify(srcUserId,isHaveRead,notifyType,classId,commId,postUserId,postCreateTime,srcContent,srcCreateTime,postContent,postUserType) values"+values

		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				if len(userIdList) != 0{
					utils.SqlTranExec(o,notifySlq,res,params)
				}
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"delete from communication where id=?",res,this.CommId)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update commComplaints set isDel = 1 where id =?",res,this.ComplaintsId)
			})
		},
	}
}


type AdminDelOneComm struct {
	AdminA
	CommIndex string
	CommId int64
	ComplaintsId string
}

func (this *AdminDelOneComm) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.ComplaintsId}
}

func (this *AdminDelOneComm)GetFun() []func(){
	var (
		commList []models.Communication
		replyContentList map[string]ReplyContent
		replyContent ReplyContent
		jsonByte []byte
		comm models.Communication
		commComplaintsList []models.CommComplaints
	)
	return []func(){
		func() {
			this.sqlQuery("select * from commComplaints where id=?",&commComplaintsList,this.ComplaintsId)
		},
		func() {
			this.checkBool = commComplaintsList[0].IsDel == 0
			this.CommId = commComplaintsList[0].CommId
			this.CommIndex = commComplaintsList[0].CommIndex
		},
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
			postContent := "管理员删除了这条信息"
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,postUserId,postCreateTime,srcUserType,postUserType,postContent) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,replyContent.PostUserId,replyContent.Content,replyContent.CreateTime,0,3,this.CommIndex,commList[0].ClassId,commList[0].Id,0,utils.GetTime(),replyContent.SrcUserType,3,postContent)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContent = ? ,replyContentTail=?,replyContentHeader=? where id = ?" ,res,string(jsonByte),comm.ReplyContentTail,comm.ReplyContentHeader,comm.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update commComplaints set isDel=1 where id=?",res,this.ComplaintsId)
			})
		},
	}
}



type AdminDelOneCommDic struct {
	AdminA
	CommIndex string
	CommId int64
}

func (this *AdminDelOneCommDic) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CommId},[]string{this.CommIndex}
}

func (this *AdminDelOneCommDic)GetFun() []func(){
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
			postContent := "管理员删除了这条信息"
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,postUserId,postCreateTime,srcUserType,postUserType,postContent) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,replyContent.PostUserId,replyContent.Content,replyContent.CreateTime,0,3,this.CommIndex,commList[0].ClassId,commList[0].Id,0,utils.GetTime(),replyContent.SrcUserType,3,postContent)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContent = ? ,replyContentTail=?,replyContentHeader=? where id = ?" ,res,string(jsonByte),comm.ReplyContentTail,comm.ReplyContentHeader,comm.Id)
			})
		},
	}
}



type AdminDeleteCommDic struct {
	AdminA
	CommId int64
}

func (this *AdminDeleteCommDic) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{}
}

func (this *AdminDeleteCommDic)GetFun() []func(){
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
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
		},
		func() {
			userIdList = append(userIdList,comm.UserId)
			for e := range replyContentList {
				//不向自己发送通知和不重复发通知
				if !utils.ValIsExistArr(userIdList, replyContentList[e].PostUserId) {
					userIdList = append(userIdList,replyContentList[e].PostUserId)
				}
			}
		},
		func() {
			length := len(userIdList)-1
			values := ""
			time := utils.GetTime()
			postContent := "管理员删除了这条信息"
			for e := range userIdList {
				if e == length{
					values += "(?,?,?,?,?,?,?,?,?,?,?)"
				}else{
					values += "(?,?,?,?,?,?,?,?,?,?,?),"
				}
				params = append(params,userIdList[e],0,3,comm.ClassId,comm.Id,0,time,comm.Content,comm.CreateTime,postContent,3)
			}
			notifySlq = "insert into replyMsgNotify(srcUserId,isHaveRead,notifyType,classId,commId,postUserId,postCreateTime,srcContent,srcCreateTime,postContent,postUserType) values"+values

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

type GetSetting struct {
	AdminA
}

func (this *GetSetting)GetFun() []func(){
	var(
		settingList []models.Setting
	)
	return []func(){
		func() {
			this.sqlQuery("select * from setting where id = 1",&settingList)
		},
		func() {
			this.setReturn(1,settingList)
		},
	}
}

type SettingModify struct {
	AdminA
	DealRatio float64
	CoinRatio float64
	SignCoin int
}
func (this *SettingModify) validArr() ([]int,[]int64,[]string) {
	return []int{this.SignCoin},[]int64{},[]string{}
}
func (this *SettingModify)GetFun() []func(){
	return []func() {
		func() {
			this.judge(this.DealRatio == 0 || this.CoinRatio == 0)
		},
		func() {
			this.sqlExec("update setting set dealRatio = ?,coinRatio=?,signCoin=?",this.DealRatio,this.CoinRatio,this.SignCoin)
		},
	}
}

type AdminAddName struct {
	AdminA
	Name string
	PicUrl string
}

func (this *AdminAddName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Name,this.PicUrl}
}

func (this *AdminAddName)GetFun() []func(){
	var(
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		userMsgJson string
	)
	return []func(){
		func() {
			this.sqlQuery("select userMsg from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			userMsg["0"] = map[string]interface{}{
				"username":this.Name,
				"img":this.PicUrl,
			}
		},
		func() {
			var byteArr []byte
			byteArr,this.checkErr = json.Marshal(userMsg)
			userMsgJson = string(byteArr)
		},
		func() {
			this.sqlExec("update `cache` set userMsg=? where id = 1",userMsgJson)
		},
	}
}

type AdminLogin struct {
	AdminA
	Username string
	Password string
}

func (this *AdminLogin) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Username,this.Password}
}
func (this *AdminLogin)GetFun() []func(){
	var(
		adminList []models.Admin
	)
	return []func(){
		func() {
			this.sqlQuery("select * from admin where name =? and password=?",&adminList,this.Username,this.Password)
		},
		func() {
			this.ctrl.SetSession("admin",adminList[0])
		},
	}
}


type AdminGetCouponList struct {
	AdminA
}

func (this *AdminGetCouponList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		couponList []models.Coupon
		start int
		countList []Count
		pageList []utils.Page
		isDo int
		isDoMap = []string{" where isShelves = 1"," where isShelves = 0",""}
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			isDo = this.ctrl.GetInt("isDo")
		},
		func() {
			this.sqlQuery("select * from `coupon` "+isDoMap[isDo]+" order by createTime desc limit ?,10 ",&couponList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `coupon`" +isDoMap[isDo],&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"couponList":couponList,
				"page":pageList,
			})
		},
	}
}

type AdminDoIsShelves struct {
	AdminA
	IsShelves int
	Id int64
}
func (this *AdminDoIsShelves) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}
func (this *AdminDoIsShelves) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("update coupon set isShelves = ? where id=?",this.IsShelves,this.Id)
		},
	}
}

type AdminAddCoupon struct {
	AdminA
	StartTime int64
	EndTime int64
	Name string
	IsShelves int8
	Amount float64
}
func (this *AdminAddCoupon) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.StartTime,this.EndTime},[]string{this.Name}
}
func (this *AdminAddCoupon) GetFun() []func() {
	return []func(){
		func() {
			this.checkBool = this.Amount != 0
		},
		func() {
			this.checkBool = this.StartTime < this.EndTime
			this.errMsg("截至日期不能大于开始日期")
		},
		func() {
			this.sqlExec("insert into coupon(name,startTime,endTime,amount,isShelves,count,createTime) values(?,?,?,?,?,?,?)",
				this.Name,this.StartTime,this.EndTime,this.Amount,this.IsShelves,0,utils.GetTime())
		},
	}
}

type GetPlatformList struct {
	AdminA
}

func (this *GetPlatformList) GetFun() []func() {
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var (
		platformNotifyList []models.PlatformNotify
		start int
		countList []Count
		pageList []utils.Page
		isDo int
		isDoMap = []string{""," where type='complaints'"," where type='withdrawal'"}
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			isDo = this.ctrl.GetInt("isDo")
		},
		func() {
			this.sqlQuery("select * from `platformNotify` "+isDoMap[isDo]+" order by createTime desc limit ?,10 ",&platformNotifyList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `platformNotify`" +isDoMap[isDo],&countList)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"platformNotifyList":platformNotifyList,
				"page":pageList,
			})
		},
	}
}

type Test struct {
	AdminA
}
func (this *Test) GetFun() []func() {
	return []func(){
		func() {
			now := utils.GetTime()
			var time int64
			//var str string
			for i := 0;i<=50;i++{
				rand := utils.RandNum(10000,99999999)
				time = now - rand

				//this.sqlExec("insert into platformNotify(title,content,createTime,type,linkId) values(?,?,?,?,?)",
				//	"投诉"+strconv.Itoa(i),"投诉",time,"complaints",i)
				//str = utils.GetRandomString(20)
				//this.sqlExec("insert into `user`(id,img,coin,createTime,username,phone,amount,password,isBan) values(?,?,?,?,?,?,?,?,?)",
				//	str,"/static/upload/user/jdOK2ZLIfB5DQmZzIHSq.jpg",0,time,"wy_"+strconv.FormatInt(rand,10),"",rand/10000,"123456",0)
				this.sqlExec("insert into teacher(name,phone,password,amount,createTime,img,unMoveAmount,sellCount,watchCount,isBan) values(?,?,?,?,?,?,?,?,?,?)",
					"讲师"+strconv.FormatInt(rand,10),"13632212810","e10adc3949ba59abbe56e057f20f883e",rand/10000,time,"/static/img/teacher/default.jpg",0,0,0,0)
			}

		},
	}
}

type AdminLogout struct {
	AdminA
}

func (this *AdminLogout) GetFun() []func() {
	return []func(){
		func() {
			this.ctrl.SetSession("admin",nil)
		},
	}
}




type AdminGetTeacherNotifyList struct {
	AdminA
}

func (this *AdminGetTeacherNotifyList) GetFun() []func() {
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
			this.sqlQuery("select * from `teacherNotify` where teacherId=0 order by createTime  desc limit ?,10 ",&teacherNotifyList,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `teacherNotify` where teacherId=0" ,&countList)
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

type AddTeacherNotify struct {
	AdminA
	Title string
	Content string
}
func (this *AddTeacherNotify) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Title,this.Content}
}
func (this *AddTeacherNotify) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("insert into teacherNotify(title,content,teacherId,createTime) values(?,?,?,?)",this.Title,this.Content,0,utils.GetTime())
		},
	}
}

type TeacherNotifyDel struct {
	AdminA
	Id int64
}

func (this *TeacherNotifyDel) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *TeacherNotifyDel) GetFun() []func() {
	return []func(){
		func() {
			this.sqlExec("delete  from teacherNotify where id=?",this.Id)
		},
	}
}
