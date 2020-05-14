package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type User1 struct {
	UserA
}

func (this *User1)GetFun() []func(){
	return []func(){
		func() {
			fmt.Println(this.user)
			fmt.Println(333)
		},
	}
}

type PostTest struct {
	UserA
	Id int64
	Name string
}

func (this *PostTest) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{this.Name}
}

func  (this *PostTest)GetFun() []func(){
	return []func(){
		func() {
			fmt.Println("33333")
			fmt.Println(this.Id,this.Name)
		},
	}

}

type CodeLogin struct {
	UserA
	Code string
}

func (this *CodeLogin) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Code}
}

func  (this *CodeLogin)GetFun() []func(){
	var(
		userList []models.User
		openId string
	)
	return []func(){
		func() {
			//fmt.Println(this.Code)
			this.checkBool,openId = utils.GetOpenId(this.Code)
		},
		func() {
			this.sqlQuery("select * from `user` where id=? ",&userList,openId)
			this.check = false
		},
		func() {
			if len(userList) > 0{
				this.ctrl.SetSession("user",userList[0])
			}else{
				img:="/static/upload/user/userDefault.jpg"
				name  := "user"+strconv.FormatInt(utils.RandNum(1000,9999),10)
				this.sqlExec("insert into `user`(id,img,coin,createTime,username,phone,amount,password,isBan) values(?,?,?,?,?,?,?,?,?)",
					openId,img,0,utils.GetTime(),name,"",0,"",0)
			}
		},
	}
}

type UserLogin struct {
	UserA
	UserId string
}
func (this *UserLogin) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.UserId}
}

type UserMsg map[string]map[string]interface{}
type UserWatchRecord map[string]bool

func  (this *UserLogin)GetFun() []func(){
	var(
		userList []models.User
		img string
		name string
		time int64
		cacheList []models.Cache
		userMsg UserMsg
		userMsgJson string
	)
	return []func(){
		func() {
			this.sqlQuery("select * from `user` where id=? ",&userList,this.UserId)
			this.check = false
		},
		func() {
			if len(userList) > 0{
				if userList[0].IsBan == 1 {
					this.setReturn(-1,"你已被管理员禁止登录")
				}else{
					this.ctrl.SetSession("user",userList[0])
					this.setReturn(1,"")
				}
			}
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			img ="/static/upload/user/userDefault.jpg"
			name  = "user"+strconv.FormatInt(utils.RandNum(1000,9999),10)
			time = utils.GetTime()
			userMsg[this.UserId] = map[string]interface{}{
				"username":name,
				"img":img,
			}
		},
		func() {
			var byteArr []byte
			byteArr,this.checkErr = json.Marshal(userMsg)
			userMsgJson = string(byteArr)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into `user`(id,img,coin,createTime,username,phone,amount,password,isBan) values(?,?,?,?,?,?,?,?,?)",
					res,this.UserId,img,0,time,name,"",0,"",0)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `cache` set userMsg = ? where id=1",res,userMsgJson)
			})
		},
		func() {
			this.ctrl.SetSession("user",models.User{this.UserId,img,0,time,name,"",0,"",0})
		},
	}
}

//func (this *CodeLogin)Next(_if func(bool), _true func([]func()),_false func([]func())){
//	_if(this.num > 0)
//	_true([]func(){
//		func() {
//			this.ctrl.SetSession("user",this.userList[0])
//		},
//	})
//	_false([]func(){
//		func() {
//			this.sqlExec("insert into `user`(id) values(?)",this.openId)
//		},
//	})
//}

type GetSession struct {
	UserA
}

func  (this *GetSession)GetFun() []func(){
	return []func(){
		func() {
			fmt.Println(this.user)
		},
	}
}

type SetUsername struct {
	UserA
	Username string
}

func (this *SetUsername) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Username}
}

func (this *SetUsername)GetFun() []func(){
	var(
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		jsons []byte
		)
	return []func(){
		func() {
			this.sqlQuery("select userMsg from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkErr  = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			userMsg[this.user.Id]["username"] = this.Username
			jsons,this.checkErr = json.Marshal(userMsg)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `user` set username = ? where id = ?",res,this.Username,this.user.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `cache` set userMsg = ? where id=1",res,string(jsons))
			})
		},
		func() {
			this.flushUser()
			this.setReturn(1,this.Username)
		},
	}
}

type UserReadPay struct {
	UserA
}

func (this *UserReadPay)GetFun() []func(){
	return []func(){
		func() {
			this.checkBool = this.user.Password != ""
		},
	}
}

type ModifyPayPwd struct {
	UserA
	Password string
}
func (this *ModifyPayPwd) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Password}
}
func (this *ModifyPayPwd)GetFun() []func(){
	return []func(){
		func() {
			this.Password = utils.Md5(this.Password)
			this.sqlExec("update `user` set password = ? where id = ?", this.Password,this.user.Id)
		},
		func() {
			this.flushUser()
		},
	}
}

type UserComplaints struct {
	UserA
	Content string
}
func (this *UserComplaints) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Content}
}
func (this *UserComplaints)GetFun() []func(){
	var(
		contentPre string
	)
	return []func(){
		func() {
			time := utils.GetTime()
			if len(this.Content) > 60{
				contentPre = utils.SubStrRunes(this.Content,20) + "..."
			} else{
				contentPre = this.Content
			}

			var complaintsId int64
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				complaintsId = utils.SqlTranExec(o,"insert into complaints(content,createTime,userId,isDo,contentPre) values(?,?,?,?,?)",res,this.Content,time,this.user.Id,0,contentPre)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into platformNotify(title,content,createTime,type,linkId) values(?,?,?,?,?)",
					res,"用户投诉，请及时处理",contentPre,time,"complaints",complaintsId)
			})
		},
	}
}

type UserGetClassByTypeId struct {
	UserA
	TypeId int64
}
func (this *UserGetClassByTypeId) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.TypeId},[]string{}
}
func (this *UserGetClassByTypeId)GetFun() []func(){
	var(
		classList []models.Class
	)
	return []func(){
		func() {
			this.sqlQuery("select * from class where typeId = ?  and isShelves=1",&classList,this.TypeId)
		},
		func() {
			this.setReturn(1,classList)
		},
	}
}

type UserClassSearch struct {
	UserA
	Content string
}

func (this *UserClassSearch) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Content}
}
func (this *UserClassSearch)GetFun() []func(){
	type Class struct {
		Name string `json:"name" orm:"column(name)"`
		Id int64 `json:"id" orm:"column(id)"`
	}
	var (
		classList []Class
	)
	return []func(){
		func() {
			this.sqlQuery(" select id,name from class where name like '%"+this.Content+"%' and isShelves = 1",&classList)
		},
		func() {
			this.setReturn(1,classList)
		},
	}
}

type VGetIndex struct {
	UserA
}

func (this *VGetIndex)GetFun() []func(){
	type Count struct {
		Value int ` orm:"column(count)"`
	}
	var(
		VPatternList []models.VPattern
		countList []Count
		pageList []utils.Page
		start int
		scopeId int64
	)

	return []func(){
		func() {
			start = this.ctrl.GetInt("start")
			scopeId,this.checkErr = this.ctrl.GetInt64("scopeId")
			if scopeId == 0{
				scopeId = 1
			}
		},
		func() {
			this.sqlQuery("select * from v_pattern where scopeId=?  order by createTime desc limit ?,10",&VPatternList,scopeId,start)
		},
		func() {
			this.sqlQuery("select count(*)  as count from `v_pattern` where scopeId=?",&countList,scopeId)
			this.check = false
			pageList = utils.GetPage(countList[0].Value,start)
		},
		func() {
			this.setReturn(1,map[string]interface{}{
				"patternList":VPatternList,
				"page":pageList,
			})
		},
	}
}


type GetVTypeList struct {
	UserA
}

func (this *GetVTypeList)GetFun() []func(){
	var(
		vTypeList []models.VType
	)
	return []func(){
		func() {
			this.sqlQuery("select * from v_type",&vTypeList)
		},
		func() {
			this.setReturn(1,vTypeList)
		},
	}
}

type VDelTypeName struct {
	UserA
	Id int64
}

func (this *VDelTypeName)GetFun() []func(){
	return []func(){
		func() {
			this.sqlIsExist("select * from v_scope where typeId=?",this.Id)
			this.errMsg("该类型下还有范围")
		},
		func() {
			this.sqlExec("delete from v_type where id=?",this.Id)
		},
	}
}

type VAddType struct {
	UserA
	Name string
}
func (this *VAddType) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{},[]string{this.Name}
}
func (this *VAddType)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("insert into v_type(name,createTime) values(?,?)",this.Name,utils.GetTime())
		},
	}
}

type VModifyName struct {
	UserA
	Name string
	Id int64
}
func (this *VModifyName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{this.Name}
}
func (this *VModifyName)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("update v_type set name = ? where id=?",this.Name,this.Id)
		},
	}
}

type VGetScopeList struct {
	UserA
}

func (this *VGetScopeList)GetFun() []func(){
	var(
		vScopeList []models.VScope
		typeId int64
	)
	return []func(){
		func() {
			typeId ,this.checkErr= this.ctrl.GetInt64("typeId")
		},
		func() {
			this.sqlQuery("select * from v_scope where typeId=?",&vScopeList,typeId)
		},
		func() {
			this.setReturn(1,vScopeList)
		},
	}
}

type VAddScope struct {
	UserA
	TypeId int64
	Name string
}
func (this *VAddScope) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.TypeId},[]string{this.Name}
}

func (this *VAddScope)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("insert into v_scope(name,typeId,createTime) values(?,?,?)",this.Name,this.TypeId,utils.GetTime())
		},
	}
}

type VDelScope struct {
	UserA
	Id int64
}
func (this *VDelScope) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}

func (this *VDelScope)GetFun() []func(){
	return []func(){
		func() {
			this.sqlIsExist("select * from v_pattern where scopeId=?",this.Id)
			this.errMsg("该范围下还有句型")
		},
		func() {
			this.sqlExec("delete from v_scope where id=?",this.Id)
		},
	}
}

type VModifyScopeName struct {
	UserA
	Id int64
	Name string
}
func (this *VModifyScopeName) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{this.Name}
}
func (this *VModifyScopeName)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("update v_scope set name = ? where id=?",this.Name,this.Id)
		},
	}
}

type VAddPattern struct {
	UserA
	Pattern string
	Example string
	ScopeId int64
}
func (this *VAddPattern) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ScopeId},[]string{this.Pattern,this.Example}
}
func (this *VAddPattern)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("insert into v_pattern(pattern,example,scopeId,createTime) values(?,?,?,?)",this.Pattern,this.Example,this.ScopeId,utils.GetTime())
		},
	}
}

type VDelPattern struct {
	UserA
	Id int64
}
func (this *VDelPattern) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.Id},[]string{}
}
func (this *VDelPattern)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("delete from v_pattern where id=?",this.Id)
		},
	}
}

type VUpdatePattern struct {
	UserA
	Pattern string
	Example string
	ScopeId int64
	Id int64
}
func (this *VUpdatePattern) validArr() ([]int,[]int64,[]string) {
	return []int{},[]int64{this.ScopeId,this.Id},[]string{this.Pattern,this.Example}
}
func (this *VUpdatePattern)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("update v_pattern set pattern = ? ,example=?, scopeId=? where id=?",this.Pattern,this.Example,this.ScopeId,this.Id)
		},
	}
}