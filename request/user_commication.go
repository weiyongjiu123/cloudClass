package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type  CommReplyOther struct {
	UserA
	CommIndex string
	CommId int64
	Content string

	commList []models.Communication
}
func (this *CommReplyOther) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.CommId},[]string{this.Content,this.CommIndex}
}
func (this *CommReplyOther)GetFun() []func(){
	var(
		comm models.Communication
		replyContentList map[string]ReplyContent
		srcReplyContent ReplyContent
		jsonByte []byte
		replyTime int64
		replyIndex string
	)
	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=? limit 1",&this.commList,this.CommId)
		},
		func() {
			comm = this.commList[0]
			replyContentStr := comm.ReplyContent
			this.checkErr = json.Unmarshal([]byte(replyContentStr),&replyContentList)
		},
		func() {
			srcReplyContent,this.checkBool = replyContentList[this.CommIndex]
		},
		func() {
			this.checkBool = srcReplyContent.PostUserId != this.user.Id
			this.baseStatus = -1
			this.baseContent = "不能回复自己"
		},
		func() {
			replyTime = utils.GetTime()
			replyIndex = utils.GetRandomString(20)
			tail := replyContentList[comm.ReplyContentTail]
			tail.Next = replyIndex
			replyContentList[comm.ReplyContentTail] = tail

			reply := ReplyContent{this.user.Id,srcReplyContent.PostUserId,this.Content,replyTime,"","","",comm.ReplyContentTail,1,srcReplyContent.PostUserType}
			replyContentList[replyIndex] = reply
			comm.ReplyContentTail = replyIndex
			jsonByte,this.checkErr = json.Marshal(replyContentList)
		},
		func() {
			//this.sqlExec("update communication set replyContentTail=?,replyContentHeader=?,replyContent=? where id=?",comm.ReplyContentTail,comm.ReplyContentHeader,string(jsonByte),comm.Id)
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContentTail=?,replyContentHeader=?,replyContent=? where id=?",res,comm.ReplyContentTail,comm.ReplyContentHeader,string(jsonByte),comm.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,postUserId,postContent,postCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,srcUserType,postUserType) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,srcReplyContent.PostUserId,srcReplyContent.Content,srcReplyContent.CreateTime,this.user.Id,this.Content,replyTime,0,1,replyIndex,comm.ClassId,comm.Id,srcReplyContent.PostUserType,1)
			})

		},
	}
}


type CommReply struct {
	UserA
	CommId int64
	Content string

	replyContentList map[string]ReplyContent
	replyContent string
	commList []models.Communication
}

func (this *CommReply) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.CommId},[]string{this.Content}
}

func (this *CommReply)GetFun() []func(){

	return []func(){
		func() {
			this.sqlQuery("select * from communication where id=? limit 1",&this.commList,this.CommId)
		},
		func() {
			this.checkBool = !(this.commList[0].UserId == this.user.Id)
			//fmt.Println(this.commList[0].UserId,"  ",this.user.Id)
			this.baseStatus = -1
			this.baseContent = "不能回复自己"
		},
	}
}
func (this *CommReply)Next(_if func(bool), _true func([]func()),_false func([]func())){
	var (
		jsons []byte
		replyIndex string
		replyTime int64
	)
	comm := this.commList[0]
	_if(comm.ReplyContentHeader != "")
	_true([]func(){
		func() {
			this.checkErr  = json.Unmarshal([]byte(comm.ReplyContent),&this.replyContentList)
		},
		func() {
			replyTime = utils.GetTime()
			reply := ReplyContent{this.user.Id,comm.UserId,this.Content,replyTime,"","","",comm.ReplyContentTail,1,comm.UserType}
			replyIndex = utils.GetRandomString(20)
			this.replyContentList[replyIndex] = reply

			replyContent := this.replyContentList[comm.ReplyContentTail]
			replyContent.Next = replyIndex
			this.replyContentList[comm.ReplyContentTail] = replyContent
			comm.ReplyContentTail = replyIndex
		},
	})
	_false([]func(){
		func() {
			this.replyContentList = map[string]ReplyContent{}
		},
		func() {
			replyTime = utils.GetTime()
			reply := ReplyContent{this.user.Id,comm.UserId,this.Content,replyTime,"","","","",1,comm.UserType}
			replyIndex = utils.GetRandomString(20)
			this.replyContentList[replyIndex] = reply

			comm.ReplyContentTail = replyIndex
			comm.ReplyContentHeader = replyIndex
		},
	})

	_if(true)
	_true([]func(){
		func() {
			jsons,this.checkErr = json.Marshal(this.replyContentList)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into replyMsgNotify(srcUserId,srcContent,srcCreateTime,postUserId,postContent,postCreateTime,isHaveRead,notifyType,replyIndex,classId,commId,srcUserType,postUserType) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
					res,comm.UserId,comm.Content,comm.CreateTime,this.user.Id,this.Content,replyTime,0,1,replyIndex,comm.ClassId,comm.Id,comm.UserType,1)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContentTail=?,replyContentHeader=?,replyContent=? where id=?",res,comm.ReplyContentTail,comm.ReplyContentHeader,string(jsons),comm.Id)
			})
		},
	})
}


type GetCommunication struct {
	UserA
}
func (this *GetCommunication)GetFun() []func(){
	var (
		classId int64
		getCommunicationRes []GetCommunicationRes
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
	)
	return []func(){
		func() {
			classId,this.checkErr = this.ctrl.GetInt64("classId")
			this.checkZero = classId
		},
		func() {
			this.sqlQuery("select replyContentHeader,replyContentTail,content,replyContent,classId,createTime,id,userId,userType from communication  where classId = ?  order by createTime DESC",&getCommunicationRes,classId)
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			cache := cacheList[0]
			if cache.UserMsg != "" && cache.UserMsg != "{}"{
				this.checkErr  = json.Unmarshal([]byte(cache.UserMsg),&userMsg)
			}else{
				//没有数据就推出
				this.checkZero = 0
			}
		},
		func() {
			var replyContentList map[string]ReplyContent
			for k,v := range getCommunicationRes {

				if v.ReplyContent != "" && v.ReplyContent != "{}"{
					er :=json.Unmarshal([]byte(v.ReplyContent),&replyContentList)
					if er != nil{
						this.checkZero = 0
					}
					for k,v := range replyContentList {
						v.PostUsername = userMsg[v.PostUserId]["username"].(string)
						v.SrcUsername = userMsg[v.SrcUserId]["username"].(string)
						replyContentList[k] = v
					}
					v.ReplyContentStrut = replyContentList
				}
				v.Username = userMsg[v.UserId]["username"].(string)
				v.Img =  userMsg[v.UserId]["img"].(string)
				getCommunicationRes[k] = v
			}
		},
		func() {
			//fmt.Println(getCommunicationRes)
			for _,v := range getCommunicationRes {
				fmt.Printf("CommId=%d   userId=%s img=%s   username=%s  留言：%s    发表时间：%d \n",v.CommId,v.UserId,v.Img,v.Username,v.Content,v.CreateTime)
				index := v.ReplyContentHeader
				for i:= 1;i>0 ; i++ {
					if index == ""{
						break
					}
					replyContent,ok := v.ReplyContentStrut[index]
					if !ok{
						break
					}
					//fmt.Println(index,"   ",replyContent)
					fmt.Printf("pre: %s next: %s      %s 回复 %s : %s    ---%s     %s\n",replyContent.Pre,replyContent.Next,replyContent.PostUsername,replyContent.SrcUsername,replyContent.Content,utils.GetDate(replyContent.CreateTime),index)
					index = replyContent.Next
				}
				//for _, v := range v.ReplyContentStrut {
				//	fmt.Printf("%s 回复 %s : %s    ---%d\n",v.PostUsername,v.SrcUsername,v.Content,v.CreateTime)
				//}
			}
		},
		func() {
			this.setReturn(1,getCommunicationRes)
		},
	}
}

//删除自己的留言，包括下面回复的也会一起删除
type DeleteCommunication struct {
	UserA
	CommId int64
}

func (this *DeleteCommunication) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.CommId},[]string{}
}
func (this *DeleteCommunication)GetFun() []func(){
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
			if comm.UserId != this.user.Id {
				this.checkZero = 0
				this.baseContent = "不可以删除别人的评论"
			}
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(comm.ReplyContent),&replyContentList)
		},
		func() {
			for e := range replyContentList {
				if this.user.Id != replyContentList[e].PostUserId &&!utils.ValIsExistArr(userIdList, replyContentList[e].PostUserId) {
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
					values += "(?,?,?,?,?,?,?,?,?,?,?)"
				}else{
					values += "(?,?,?,?,?,?,?,?,?,?,?),"
				}
				params = append(params,userIdList[e],0,2,comm.ClassId,comm.Id,this.user.Id,time,comm.Content,comm.CreateTime,comm.UserType,1)
			}
			notifySlq = "insert into replyMsgNotify(srcUserId,isHaveRead,notifyType,classId,commId,postUserId,postCreateTime,srcContent,srcCreateTime,srcUserType,postUserType) values"+values

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

type GetReplyNotify struct {
	UserA
}
func (this *GetReplyNotify)GetFun() []func(){
	var (
		replyMsgNotifyList []models.ReplyMsgNotify
		cacheList []models.Cache
		userMsgArr map[string]map[string]interface{}
		)
	return []func(){
		func() {
			this.sqlQuery("select * from replyMsgNotify where srcUserId = ? and isHaveRead=0 order by postCreateTime desc",&replyMsgNotifyList,this.user.Id)
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

type GetMsgDetail struct {
	UserA
}


func (this *GetMsgDetail)GetFun() []func(){
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


type AddCommunication struct {
	UserA
	Content string
	ClassId int64
}


func (this *AddCommunication) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ClassId},[]string{this.Content}
}

func (this *AddCommunication)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("insert into communication(userId,content,replyContent,classId,createTime,userType) values(?,?,?,?,?,?)",this.user.Id,this.Content,"{}",this.ClassId,utils.GetTime(),1)
		},
	}
}

type GetPerReplyNotify struct {
	UserA
	cacheList []models.Cache
	replyMsgNotifyList []models.ReplyMsgNotify
}

func (this *GetPerReplyNotify)GetFun() []func(){
	var (


	)
	return []func(){
		func() {
			this.sqlQuery("select * from replyMsgNotify where srcUserId = ?",&this.replyMsgNotifyList,this.user.Id)
		},
		func() {
			this.sqlQuery("select * from cache where id=1",&this.cacheList)
		},

	}
}

func (this *GetPerReplyNotify)Next(_if func(bool), _true func([]func()),_false func([]func())){
	var (
		userMsgMap map[string]map[string]interface{}
		convBool bool
	)
	_if(this.cacheList[0].UserMsg == "")
	_true([]func(){
		//直接返回
		func() {
			for _, v := range this.replyMsgNotifyList {
				fmt.Println(v)
			}
		},
	})
	_false([]func(){
		func() {
			this.checkErr = json.Unmarshal([]byte(this.cacheList[0].UserMsg),&userMsgMap)
		},
		func() {
			for k, v := range this.replyMsgNotifyList {
				v.SrcUsername,convBool = userMsgMap[v.SrcUserId]["username"].(string)
				v.SrcUserImg,convBool = userMsgMap[v.SrcUserId]["img"].(string)
				v.PostUsername,convBool = userMsgMap[v.PostUserId]["username"].(string)
				v.PostUserImg,convBool = userMsgMap[v.PostUserId]["img"].(string)
				this.replyMsgNotifyList[k] = v
			}
		},
		func() {
			for _, v := range this.replyMsgNotifyList {
				fmt.Println(v)
			}
		},
	})
}


type TestSetNameImg struct {
	UserA
}

func (this *TestSetNameImg)GetFun() []func(){
	var userMsgByte []byte
	return []func(){
		func() {
			userMsg := map[string]map[string]interface{}{
				"testOpenId":{
					"username":"新名字33",
					"img":"img_10",
				},
				"oeg225I9auRLzcI4QMXELJdVYKJc":{
				"username":"小红",
				"img":"img_11",
			},
			}
			userMsgByte,this.checkErr = json.Marshal(userMsg)
		},
		func() {
			this.sqlExec("update cache set userMsg = ? where id = 1",string(userMsgByte))
		},
	}
}

//删除自己回复他人的留言
type CommDelete struct {
	UserA
	CommIndex string
	CommId int64
}
func (this *CommDelete) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CommId},[]string{this.CommIndex}
}
func (this *CommDelete)GetFun() []func(){
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
			this.checkBool = replyContent.PostUserId == this.user.Id
			this.baseContent = "不能删除别人的留言"
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
					res,replyContent.SrcUserId,replyContent.Content,replyContent.CreateTime,0,2,this.CommIndex,commList[0].ClassId,commList[0].Id,this.user.Id,utils.GetTime(),replyContent.PostUserType,1)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update communication set replyContent = ? ,replyContentTail=?,replyContentHeader=? where id = ?" ,res,string(jsonByte),comm.ReplyContentTail,comm.ReplyContentHeader,comm.Id)
			})
		},
	}
}



type UserCommPlaints struct {
	UserA
	CommId int64
	CommIndex string
}

func (this *UserCommPlaints) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CommId},[]string{this.CommIndex}
}

func (this *UserCommPlaints)GetFun() []func(){
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
			if len(allContent)>60 {
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


