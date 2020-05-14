package request

import (
	"cloudClass/models"
	"cloudClass/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type GetHomeData struct {
	UserA
}

type GetHomeDataTypeRes struct {
	Id int64 `json:"typeId" orm:"column(id)"`
	Name string `json:"text" orm:"column(name)"`
}
type GetHomeDataClassRes struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	Price float64 `json:"price" orm:"column(price)"`
	IsCharge int8  `json:"isCharge" orm:"column(isCharge)"`
	ClassImg string  `json:"classImg" orm:"column(classImg)"`
}

func (this *GetHomeData)GetFun() []func()   {
	var (
		getHomeDataTypeResList []GetHomeDataTypeRes
		getHomeDataClassResList []GetHomeDataClassRes
		adList []models.Ad
	)
	return []func(){
		func() {
			this.sqlQuery("select id,name from `type`",&getHomeDataTypeResList)
		},
		func() {
			this.sqlQuery("select id,name,price,isCharge,classImg from `class` where isShelves=1 order by watchCount desc limit 0 ,6",&getHomeDataClassResList)
		},
		func() {
			this.sqlQuery("select * from ad ",&adList)
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"typeList":getHomeDataTypeResList,
				"classList":getHomeDataClassResList,
				"adList":adList,
			})
		},
	}
}

type GetClassData struct {
	UserA
}

type GetClassDataClassList struct {
	Id int64  `json:"id" orm:"column(id)"`
	ClassImg string `json:"classImg" orm:"column(classImg)"`
	Name string `json:"name" orm:"column(name)"`
	Detail string `json:"detail" orm:"column(detail)"`
	Star float32	 `json:"star" orm:"column(star)"`
	IsCharge int8  `json:"isCharge" orm:"column(isCharge)"`
	Price float32 `json:"price" orm:"column(price)"`

}

func (this *GetClassData)GetFun() []func() {
	var (
		classList []GetClassDataClassList
		classId int64
		orderList []models.Order
	)
	return []func(){
		func() {
			classId,this.checkErr  = this.ctrl.GetInt64("classId")
			this.checkZero = classId
		},
		func() {
			this.sqlQuery("select id,classImg,name,detail,star,isCharge,price from `class` where id=? and isShelves=1 limit 1",&classList,classId)
		},
		func() {
			this.sqlQuery("SELECT * FROM `order` WHERE	userId = ? AND classId = ?  limit 1",&orderList,this.user.Id,classId)
			this.check = false
		},
		func() {
			var status int8
			if len(orderList) > 0{
				status = orderList[0].Status
			}
			this.setReturn(1, map[string]interface{}{
				"classMsg":classList[0],
				"status":status,
			})
		},
	}
}

type GetCatalog struct {
	UserA
}

func (this *GetCatalog)GetFun() []func()   {
	type CatalogList struct {
		ChapterTitleId int64  `json:"chapterTitleId" orm:"column(chapterTitleId)"`
		Name string  `json:"name" orm:"column(name)"`
		ChapterList string `json:"chapterList" orm:"column(chapterList)"`
	}
	type Class struct {
		UserWatchRecord string  `json:"userWatchRecord" orm:"column(userWatchRecord)"`
	}

	var(
		catalogListList []CatalogList
		classId int64
		classList []Class
	)
	return []func(){
		func() {
			classId,this.checkErr = this.ctrl.GetInt64("classId")
			this.checkZero = classId
		},
		func() {
			this.sqlQuery("select userWatchRecord from class  where id=?",&classList,classId)
		},
		func() {
			this.sqlQuery("select t.id as chapterTitleId,t.`name` ,(select GROUP_CONCAT(concat_ws('~_',c.id,c.name,c.time)) from chapter as c where c.chapterTitleId= t.id) as chapterList from chaptertitle as t where t.classId = ? ",&catalogListList,classId)
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"catalogList":catalogListList,
				"userWatchRecord":classList[0].UserWatchRecord,
			})
		},
	}
}

type AddEvaluation struct {
	UserA
	Star int64
	Content string
	OrderId string
}

func (this *AddEvaluation) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.Star},[]string{this.Content,this.OrderId}
}

func (this *AddEvaluation)GetFun() []func(){
	type Order struct{
		ClassId int64 `orm:"column(classId)"`
		Status int8 ` orm:"column(status)"`
	}
	type Class struct {
		EvaluationCount int64 `json:"evaluationCount" orm:"column(evaluationCount)"`
		StarSum int64 `json:"starSum" orm:"column(starSum)"`
	}
	var(
		orderList []Order
		classList []Class
	)
	return []func(){
		func() {
			this.sqlQuery("select classId,status from `order` where orderId =?",&orderList,this.OrderId)
		},
		func() {

			if orderList[0].Status != 2{
				this.checkZero = 0
			}
		},
		func() {
			this.sqlQuery("select evaluationCount,starSum from class where id=? ",&classList,orderList[0].ClassId)
		},
		func() {
			evaluationCount := classList[0].EvaluationCount + 1
			starSum := classList[0].StarSum + this.Star
			star :=  float64(starSum) / float64(evaluationCount)
			fmt.Println(starSum,evaluationCount,star)
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into evaluation(userId,star,content,classId,orderId,createTime) values(?,?,?,?,?,?)",res,this.user.Id,this.Star,this.Content,orderList[0].ClassId,this.OrderId,utils.GetTime())
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update class set evaluationCount = ? ,starSum = ?,star=? where id = ?",res,evaluationCount,starSum,star,orderList[0].ClassId)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `order` set status=3 where orderId =?",res,this.OrderId)
			})
		},
	}
}

type VideoPlayer struct {
	UserA
	VideoUrl string
}

func (this *VideoPlayer)GetFun() []func() {
	var(
		cacheList []Cache
		videoPerMap map[string]VideoPer
		videoPer VideoPer
		videoPerMapByte []byte
	)
	return []func(){
		func() {
			this.VideoUrl = this.ctrl.GetString("videoUrl")
			this.checkString = this.VideoUrl
		},
		func() {
			this.sqlQuery("select videoPer from `cache` where id=1",&cacheList)
		},
		func() {
			this.checkString = cacheList[0].VideoPer
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].VideoPer),&videoPerMap)
		},
		func() {
			videoPer,this.checkBool = videoPerMap[this.VideoUrl]
		},
		func() {
			timeLimit := videoPer.Time + 7200 //两个小时
			//this.checkBool = videoPer.Time < timeLimit
			this.checkBool = utils.GetTime() < timeLimit
			this.baseStatus = -1
			this.baseContent = "请重新进入"
		},
		func() {
			videoPer.Time  = videoPer.Time + 7200
			videoPerMap[this.VideoUrl] = videoPer
		},
		func() {
			videoPerMapByte,this.checkErr = json.Marshal(videoPerMap)
		},
		func() {
			this.sqlExec("update `cache` set videoPer = ? where id = 1",string(videoPerMapByte))
		},
		func() {
			this.setVideoUrl(videoPer.Url)
		},
	}
}

type AddStudyRecord struct {
	UserA
	ChapterId int64
	Title string
}

func (this *AddStudyRecord) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.ChapterId},[]string{this.Title}
}

func (this *AddStudyRecord)GetFun() []func(){
	type Chapter struct {
		ClassId int64 `json:"classId" orm:"column(classId)"`
		Video string  `json:"video" orm:"column(video)"`
		TeacherId int64  `json:"teacherId" orm:"column(teacherId)"`
	}
	type UserClass struct {
		UserId string `json:"userId" orm:"column(userId)"`
		ClassId int64  `json:"classId" orm:"column(classId)"`
	}


	var(
		chapterList []Chapter
		orderList []models.Order
		cacheList []Cache
		videoPerMag map[string]VideoPer
		videoIndexRandom string
		videoPerByte []byte
		classList []models.Class
		userWatchRecord UserWatchRecord
		userWatchRecordStr string
	)
	return []func(){
		func() {
			this.sqlQuery("select classId,video,teacherId from chapter where id = ? limit 1",&chapterList,this.ChapterId)
		},
		func() {
			this.sqlQuery("select * from class where id=?",&classList,chapterList[0].ClassId)
		},
		func() {
			//this.sqlQuery("select * from userClass where userId = ? and classId = ?",&userClassList,this.user.Id,chapterList[0].ClassId)
			if classList[0].IsCharge == 1 {
				this.sqlQuery("select * from `order` where userId=? and classId =?",&orderList,this.user.Id,chapterList[0].ClassId)
				this.baseStatus = -1
				this.baseContent = "先购买课程再观看"
			}

		},
		func() {
			if classList[0].UserWatchRecord == "" {
				classList[0].UserWatchRecord ="{}"
			}
			this.checkErr = json.Unmarshal([]byte(classList[0].UserWatchRecord),&userWatchRecord)
		},
		func() {
			var tmpByte []byte
			index := this.user.Id +"-" + strconv.FormatInt(classList[0].Id,10) + "-" + strconv.FormatInt(this.ChapterId,10)
			userWatchRecord[index] = true
			tmpByte,this.checkErr = json.Marshal(userWatchRecord)
			userWatchRecordStr = string(tmpByte)
		},
		func() {
			this.sqlQuery("select videoPer from `cache` where id = 1",&cacheList)
		},
		func() {
			if cacheList[0].VideoPer == ""{
				videoPerMag = map[string]VideoPer{}
			}else{
				this.checkErr = json.Unmarshal([]byte(cacheList[0].VideoPer),&videoPerMag)
			}
			videoIndexRandom = utils.GetRandomString(20)
			videoPerMag[videoIndexRandom] = VideoPer{utils.GetTime(),chapterList[0].Video}
			videoPerByte,this.checkErr = json.Marshal(videoPerMag)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into studyRecord(classId,chapterId,createTime,title,userId) values(?,?,?,?,?)",res,chapterList[0].ClassId,this.ChapterId,utils.GetTime(),this.Title,this.user.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `cache` set videoPer=? where id=1",res,string(videoPerByte))
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update chapter set watchCount = watchCount + 1 where id=?",res,this.ChapterId)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update teacher set watchCount = watchCount + 1 where id = ?",res,chapterList[0].TeacherId)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update class set watchCount = watchCount+1,userWatchRecord =? where id=?",res,userWatchRecordStr,chapterList[0].ClassId)
			}, func(o orm.Ormer, res *bool) {
				date := utils.GetNowDateA()
				time := utils.GetTime()
				_,_,affectRow := utils.SqlExec("update watchRecord set watchCount = watchCount + 1 where teacherId = ? and date=?",chapterList[0].TeacherId,date)
				if affectRow == 0 {
					utils.SqlTranExec(o,"insert into watchRecord(date,watchCount,createTime,teacherId) values(?,?,?,?)",res,date,1,time,chapterList[0].TeacherId)
				}
				fmt.Println("affectRow ",affectRow,date)
				//utils.SqlExec(")
			})
		},
		func() {
			this.setReturn(1,videoIndexRandom)
		},
	}
}

type GetStudyRecord struct {
	UserA
}

func (this *GetStudyRecord)GetFun() []func(){
	type StudyRecord struct{
		ClassId int64 `json:"classId" orm:"column(classId)"`
		ChapterId int64  `json:"chapterId" orm:"column(chapterId)"`
		CreateTime int64 `json:"createTime" orm:"column(createTime)"`
		ClassImg string `json:"classImg" orm:"column(classImg)"`
		Title string `json:"title" orm:"column(title)"`
		Name string `json:"className" orm:"column(name)"`
	}
	var studyRecordList   []StudyRecord
	return []func(){
		func() {
			this.sqlQuery("select c.name,r.createTime,r.classId,r.chapterId,r.title,c.classImg from studyRecord as r  inner join `class` as c on r.classId = c.id and r.userId = ? order by r.createTime desc limit 13",&studyRecordList,this.user.Id)
		},
		func() {
			this.setReturn(1,studyRecordList)
		},
	}
}

type AddShopCar struct {
	UserA
	ClassId int64
}

func (this *AddShopCar) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.ClassId},[]string{}
}

func (this *AddShopCar)GetFun() []func(){
	var shopCarList []models.ShopCar
	return []func(){
		func() {
			this.sqlIsExist("select * from shopCar where classId=? and userId = ? limit 1",&shopCarList,this.ClassId,this.user.Id)
			this.baseStatus = -1
			this.baseContent = "你已经添加到购物车了"
		},
		func() {
			this.sqlExec("insert into shopCar(classId,userId,createTime) values(?,?,?)",this.ClassId,this.user.Id,utils.GetTime())
		},
	}
}

type GetShopCar struct {
	UserA
}

func (this *GetShopCar)GetFun() []func(){
	type ShopCarRes struct {
		Id int64  `json:"id" orm:"column(id)"`
		UserId string `json:"-" orm:"column(userId)"`
		ClassId int64 `json:"classId" orm:"column(classId)"`
		CreateTime int64  `json:"createTime" orm:"column(createTime)"`
		Name string `json:"name" orm:"column(name)"`
		ClassImg string `json:"classImg" orm:"column(classImg)"`
		Price float64 `json:"price" orm:"column(price)"`
	}
	var shopCarList []ShopCarRes
	return []func(){
		func() {
			this.sqlQuery("select s.*,c.name,c.classImg,c.price from shopCar as s inner join class as c on s.userId = ? and s.classId = c.id",&shopCarList,this.user.Id)
		},
		func() {
			this.setReturn(1,shopCarList)
		},
	}
}

type ShopCarDelete struct {
	UserA
	ShopCarId int64
}

func (this *ShopCarDelete) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.ShopCarId},[]string{}
}

func (this *ShopCarDelete)GetFun() []func(){
	return []func(){
		func() {
			this.sqlExec("delete  from shopCar where id = ? and userId = ?",this.ShopCarId,this.user.Id)
		},
	}
}

//提交订单并支付，有两种提交
type AddOrder struct {
	UserA
	ClassIdList []int64
	Password string
}
func (this *AddOrder) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{},[]string{this.Password}
}
func (this *AddOrder)GetFun() []func(){
	type Class struct {
		Name string `json:"name"`
		ClassImg string `json:"classImg" orm:"column(classImg)"`
		Price float64 `json:"price"`
		Id int64 `json:"id"`
		TeacherId int64 `json:"teacherId" orm:"column(teacherId)"`
	}
	type Order struct {
		ClassId int64 `orm:"column(classId)"`
	}

	type OrderCoupon struct {
		CouponIdList []int64
		IsUpdate bool
	}
	type Bill struct {
		UserId string
		Amount float64
		CouponList []models.UserCoupon
		ClassList []Class
		Coin float64
		CouponListJsonStr string
		ClassListJsonStr string
		Pay float64
		Id string
		orderIdList []string
		orderIdListJsonStr string
	}
	type OrderUser struct {
		Coin float64
		IsUpdate bool
		Amount float64
	}

	type Teacher struct {
		Id int64
		Amount float64
	}

	var(
		length int
		classList []Class
		orderList []Order
		userCouponList []models.UserCoupon
		allAmount float64
		settingList []models.Setting
		couponListSql string
		couponIdListLen int
		orderCoupon OrderCoupon
		orderUser OrderUser
		bill Bill
		teacherList []Teacher
		argsMap [][]interface{}
		teacherBillParams  [][]interface{}
		time int64
	)
	return []func(){
		func() {
			//检查classId是否存在为0的情况
			this.validInt64Arr(this.ClassIdList)
		},
		func() {
			this.checkBool = this.user.Password  == utils.Md5(this.Password)
			this.errMsg("密码错误")
		},
		func() {
			//检测所提交的课程有没有在 未评价 已完成中的订单
			this.sqlQuery("select classId from `order` where userId = ? ",&orderList,this.user.Id)
			this.check = false
		},
		func() {
			fmt.Println(orderList)
			//检查是否已经购买
			flag := true
			for e := range orderList {
				for i := range this.ClassIdList {
					if orderList[e].ClassId == this.ClassIdList[i] {
						flag = false
						break
					}
				}
			}
			this.checkBool = flag
			this.errMsg("您已经购买了，无需要再购买")
		},
		func() {
			//获取课程的信息
			str := "("
			length = len(this.ClassIdList)
			classIdListLen := length - 1
			for e := range this.ClassIdList {
				if e <classIdListLen{
					str += strconv.FormatInt(this.ClassIdList[e],10) + ","
				}else{
					str += strconv.FormatInt(this.ClassIdList[e],10) + ")"
				}
			}
			this.sqlQuery("select id,name,classImg,price,teacherId from class where id in "+str,&classList)
		},
		func() {
			//检查所有的课程id是否都存在，有一个不存在都不能往下执行
			this.judge(length != len(classList))
		},
		func() {
			now := utils.GetTime()
			this.sqlQuery("SELECT * FROM userCoupon WHERE userId = ? AND isUsered = 0 AND  ? BETWEEN startTime AND endTime ",&userCouponList,this.user.Id,now)
			this.check = false
		},
		func() {
			this.sqlQuery("SELECT * FROM setting WHERE id = 1",&settingList)
		},
		func() {
			bill.ClassList = classList
			allAmount = 0
			for e := range classList {
				allAmount += classList[e].Price
				teacherList = append(teacherList,Teacher{classList[e].TeacherId,classList[e].Price})
			}
			bill.Amount = allAmount
			fmt.Println("总金额：",allAmount)
			//优惠卷
			for e := range userCouponList {
				if allAmount > userCouponList[e].Amount {
					allAmount -= userCouponList[e].Amount
					orderCoupon.IsUpdate = true
					orderCoupon.CouponIdList = append(orderCoupon.CouponIdList,userCouponList[e].Id)
					bill.CouponList = append(bill.CouponList,userCouponList[e])
				}

			}
			fmt.Println("扣除优惠劵后：",allAmount)
			//积分
			if this.user.Coin > 0 && allAmount > 0{
				coinAmount := this.user.Coin * settingList[0].CoinRatio
				if coinAmount <= allAmount{
					allAmount -= coinAmount
					orderUser.IsUpdate = true
					orderUser.Coin = 0
					bill.Coin = this.user.Coin
				}else {
					orderUser.Coin  = (coinAmount - allAmount) /  settingList[0].CoinRatio
					bill.Coin = allAmount / settingList[0].CoinRatio
					allAmount = 0
					orderUser.IsUpdate = true
				}
			}
			fmt.Println("扣除积分后：",allAmount)

			if allAmount != 0 &&this.user.Amount < allAmount {
				this.checkZero = 0
				this.errMsg("您的余额不足")
			}else if allAmount == 0{
				orderUser.Amount = this.user.Amount
			}else if allAmount != 0{
				bill.Pay = allAmount
				orderUser.Amount = this.user.Amount - allAmount
				orderUser.IsUpdate = true
			}
			bill.UserId = this.user.Id
			fmt.Println("用户金额：",orderUser.Amount)
			fmt.Println("bill : ",bill)
		},
		func() {
			var (
				classListByte []byte
			)
			classListByte ,this.checkErr =json.Marshal(bill.ClassList)
			bill.ClassListJsonStr = string(classListByte)
			bill.Id = strconv.FormatInt(utils.GetTime(),10) + strconv.FormatInt(utils.RandNum(1000000000,9999999999),10)
		},
		func() {
			if !orderCoupon.IsUpdate {
				bill.CouponListJsonStr = "[]"
				return
			}
			var couponListByte []byte
			couponListByte ,this.checkErr = json.Marshal(bill.CouponList)
			bill.CouponListJsonStr = string(couponListByte)
		},
		func() {
			couponSqlWhere := "("
			couponIdListLen = len(orderCoupon.CouponIdList) -1

			for e := range orderCoupon.CouponIdList {
				if couponIdListLen == e {
					couponSqlWhere += "?)"
				}else{
					couponSqlWhere += "?,"
				}
			}
			if orderCoupon.IsUpdate {
				couponListSql = "UPDATE userCoupon SET isUsered = 1 WHERE id in " + couponSqlWhere
			}
			//fmt.Println(couponListSql)
			//this.setReturn(1,"")
		},
		func() {
			for _, v := range teacherList {
				this.sqlExec("update teacher set amount = amount + ?,sellCount=sellCount+1 where id = ?",v.Amount * settingList[0].DealRatio,v.Id)
			}
		},
		func() {
			time = utils.GetTime()
			for _, v := range classList {
				orderId := strconv.FormatInt(utils.GetTime(),10) + strconv.FormatInt(utils.RandNum(100000000,999999999),10)
				bill.orderIdList = append(bill.orderIdList,orderId)
				item := []interface{}{1,2,time,v.Price,v.Id,v.Name,v.ClassImg,this.user.Id,orderId,v.TeacherId,bill.Id}
				argsMap = append(argsMap,item)
				amount := v.Price * settingList[0].DealRatio
				billItem := []interface{}{"《"+v.Name+"》课程出售1份,课程价格为￥"+utils.FloatToString(v.Price)+",收益比率为"+utils.FloatToString(settingList[0].DealRatio)+",实际收入为￥"+utils.FloatToString(amount),amount,time,orderId,v.TeacherId}
				teacherBillParams = append(teacherBillParams,billItem)
			}
		},
		func() {
			var byteArr []byte
			byteArr,this.checkErr = json.Marshal(bill.orderIdList)
			bill.orderIdListJsonStr = string(byteArr)
		},
		func() {

			num := length-2
			sql := "INSERT INTO `order` (isPay,`status`,createTime,amount,classId,className,classImg,userId,orderId,teacherId,billId) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
			teacherBillSql := "insert into teacherBill(description,amount,createTime,orderId,teacherId) values(?,?,?,?,?)"
			for i:=0;i<= num;i++{
				sql += ",(?,?,?,?,?,?,?,?,?,?,?)"
				teacherBillSql += ",(?,?,?,?,?)"
			}
			carSql := "delete from shopCar where  userId=? and  classId in ("
			carParam := []interface{}{}
			carLent := len(this.ClassIdList) -1
			carParam =  append(carParam, this.user.Id)
			for e := range this.ClassIdList {
				if carLent == e{
					carSql += "?)"
				}else{
					carSql += "?,"
				}
				carParam = append(carParam, this.ClassIdList[e])
			}


			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,sql,res,argsMap)
			}, func(o orm.Ormer, res *bool) {
				if orderCoupon.IsUpdate{
					utils.SqlTranExec(o,couponListSql,res,orderCoupon.CouponIdList)
				}
			}, func(o orm.Ormer, res *bool) {
				if orderUser.IsUpdate{
					utils.SqlTranExec(o,"UPDATE `user` set amount = ?,coin = ? WHERE id = ? ",res,orderUser.Amount,orderUser.Coin,this.user.Id)
				}
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into bill(id,userId,amount,couponList,classList,coin,pay,coinRatio,orderIdList,createTime) values(?,?,?,?,?,?,?,?,?,?)",res,bill.Id,bill.UserId,bill.Amount,bill.CouponListJsonStr,bill.ClassListJsonStr,bill.Coin,bill.Pay,settingList[0].CoinRatio,bill.orderIdListJsonStr,time)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,teacherBillSql,res,teacherBillParams)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,carSql,res,carParam)
				*res = true
			})
			//this.sqlExec(sql,argsMap) //批处理 插入多条数据
		},
		func() {
			this.flushUser()
		},
	}
}

type ReadPay struct {
	UserA
}

func (this *ReadPay)GetFun() []func(){
	type Coupon struct {
		Name string `json:"name" orm:"column(name)"`
		Amount float64 `json:"amount" orm:"column(amount)"`
	}
	var (
		userCouponList []models.UserCoupon
		couponList []Coupon
		settingList []models.Setting
	)
	return []func(){
		func() {
			this.sqlQuery("select * from userCoupon where userId = ? and isUsered=0",&userCouponList,this.user.Id)
			this.check = false
		},
		func() {
			this.sqlQuery("select * from setting where id = 1",&settingList)
		},
		func() {
			nowTime := utils.GetTime()
			if len(userCouponList) != 0{
				for e := range userCouponList {
					if nowTime < userCouponList[e].EndTime {
						couponList = append(couponList,Coupon{userCouponList[e].Name,userCouponList[e].Amount})
					}
				}
			}
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"couponList":couponList,
				"coin":this.user.Coin,
				"coinRatio":settingList[0].CoinRatio,
			})
		},
	}
}

type UserAddCoupon struct {
	UserA
	CouponId int64
}
func (this *UserAddCoupon) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.CouponId},[]string{}
}

func (this *UserAddCoupon)GetFun() []func(){
	var (
		couponList []models.Coupon
	)
	return []func(){
		func() {
			this.sqlQuery("SELECT * FROM coupon where id = ? LIMIT 1",&couponList,this.CouponId)
			this.baseStatus = -1
			this.baseContent = "领取失败，请稍后再试"
		},
		func() {
			this.checkZero = int64(couponList[0].IsShelves)
			this.baseStatus = -1
			this.baseContent = "该优惠卷已经下架"
		},
		func() {
			now := utils.GetTime()
			this.checkBool = now <= couponList[0].EndTime
			this.baseStatus = -1
			this.baseContent = "该优惠卷已经过期，不能领取"
		},
		func() {
			this.sqlIsExist("SELECT * FROM userCoupon WHERE couponId = ? AND userId = ?",this.CouponId,this.user.Id)
			this.baseStatus = -1
			this.baseContent = "您已领取了该优惠卷"
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"INSERT INTO userCoupon (name,startTime,endTime,amount,userId,couponId,isUsered) VALUES(?,?,?,?,?,?,?)",res,
					couponList[0].Name,couponList[0].StartTime,couponList[0].EndTime,couponList[0].Amount,this.user.Id,this.CouponId,0)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update coupon set count=count+1 where id=?",res,this.CouponId)
			})

		},
	}
}

type UserPay struct {
	UserA
	OrderId int64
}

func (this *UserPay) validArr() ([]int,[]int64,[]string){
	return []int{},[]int64{this.OrderId},[]string{}
}

func (this *UserPay)GetFun() []func(){
	var(
		orderList []models.Order
		isUseredList []models.UserCoupon
		amount float64
		settingList []models.Setting
	)
	return []func(){
		func() {
			this.sqlQuery("SELECT * FROM class where id = ? LIMIT 1",&orderList,this.OrderId)
		},
		func() {
			this.sqlQuery("SELECT * FROM userCoupon WHERE userId=? and isUsered=0",&isUseredList,this.user.Id)
			this.check = false
		},
		func() {
			this.sqlQuery("SELECT * FROM setting where id =1",&settingList)
		},
		func() {
			amount = orderList[0].Amount
			for e := range isUseredList {
				amount -= isUseredList[e].Amount
			}
			amount -= this.user.Coin/10

			fmt.Println(amount)
			fmt.Println(this.user.Amount)
			fmt.Println(settingList)
			fmt.Println(orderList)
		},
	}
}

type GetBill struct {
	UserA
}

func (this *GetBill)GetFun() []func(){
	var billList []models.Bill
	return []func(){
		func() {
			this.sqlQuery("select * from bill where userId = ? order by createTime desc",&billList,this.user.Id)
		},
		func() {
			this.setReturn(1,billList)
		},
	}
}

type GetOrder struct {
	UserA
}

func (this *GetOrder)GetFun() []func(){
	type Order struct {
		Id string `json:"id" orm:"column(orderId)"`
		Status int `json:"status" orm:"column(status)"`
		CreateTime int64 `json:"createTime" orm:"column(createTime)"`
		Amount float64 `json:"amount" orm:"column(amount)"`
		ClassName string `json:"className" orm:"column(className)"`
		ClassImg string `json:"classImg" orm:"column(classImg)"`
		ClassId int64   `json:"classId" orm:"column(classId)"`
	}
	var (
		orderList []Order
		orderType int
	)
	return []func(){
		func() {
			orderType = this.ctrl.GetInt("orderType")
			if orderType == 0{
				this.sqlQuery("select orderId,status,createTime,amount,className,classImg,classId from `order` where userId = ? order by createTime desc",&orderList,this.user.Id)
			}else if orderType == 1{
				this.sqlQuery("select * from `order` where userId = ? and status = 2 order by createTime desc",&orderList,this.user.Id)
			}else if orderType == 2{
				this.sqlQuery("select * from `order` where userId = ? and status = 3 order by createTime desc",&orderList,this.user.Id)
			}else {
				this.checkZero = 0
			}
		},
		func() {
			this.setReturn(1,orderList)
		},
		func() {
			fmt.Println(orderList)
		},
	}
}

type GetOneOrder struct {
	UserA
}

func (this *GetOneOrder)GetFun() []func(){
	var(
		orderList []models.Order
		orderId string
	)
	return []func(){
		func() {
			orderId = this.ctrl.GetString("orderId")
			this.checkString = orderId
		},
		func() {
			this.sqlQuery("select * from `order` where orderId = ?",&orderList,orderId)
		},
		func() {
			this.setReturn(1, map[string]interface{}{
				"className":orderList[0].ClassName,
				"classImg":orderList[0].ClassImg,
			})
		},
	}
}

type GetEvaluationList struct {
	UserA
}

func (this *GetEvaluationList)GetFun() []func(){
	type Evaluation struct {
		Id int64 `json:"id" orm:"column(id)"`
		UserId string `json:"userId" orm:"column(userId)"`
		Star int8 `json:"star" orm:"column(star)"`
		Content string `json:"content" orm:"column(content)"`
		OrderId string `json:"-" orm:"column(orderId)"`
		Username string `json:"username" `
		UserImg string `json:"userImg" `
		CreateTime int64  `json:"createTime" orm:"column(createTime)"`
	}

	var(
		classId int64
		evaluationList[] Evaluation
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
	)
	return []func(){
		func() {
			classId,this.checkErr = this.ctrl.GetInt64("classId")
			this.checkZero = classId
		},

		func() {
			this.sqlQuery("select * from evaluation where classId = ?",&evaluationList,classId)
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			for k, v := range evaluationList {
				v.Username,this.checkBool = userMsg[v.UserId]["username"].(string)
				v.UserImg,this.checkBool = userMsg[v.UserId]["img"].(string)
				evaluationList[k] = v
			}
			fmt.Println(userMsg)
			this.check = false
		},
		func() {

			this.setReturn(1,evaluationList)
		},
	}
}

type Sign struct {
	UserA
}

func (this *Sign)GetFun() []func(){
	var(
		date string
		time int64
		settingList []models.Setting
	)
	return []func(){
		func() {
			time = utils.GetTime()
			date = utils.GetNowDate()
			this.sqlIsExist("select * from `sign` where userId = ? and date = ? limit 1",this.user.Id,date)
			this.baseContent = "你今日已经签到了"
			this.baseStatus = -1
		},
		func() {
			this.sqlQuery("select signCoin from setting where id = 1",&settingList)
		},
		func() {
			fmt.Println("coin: ",this.user.Coin)
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"insert into `sign`(userId,date,createTime,coin) values(?,?,?,?)",res,this.user.Id,date,time,settingList[0].SignCoin)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `user` set coin=coin + ? where id=?",res,settingList[0].SignCoin,this.user.Id)
			})
		},
		func() {
			this.flushUser()
		},
	}
}

type GetSignList struct {
	UserA
}

func (this *GetSignList)GetFun() []func(){
	type Sign struct {
		Date string `json:"date"`
		Coin int64 `json:"coin"`
	}
	var(
		signList []Sign
	)
	return []func(){
		func() {
			this.sqlQuery("select date,coin from `sign` where userId= ? order by createTime desc",&signList,this.user.Id)
		},
		func() {
			this.setReturn(1,signList)
		},
	}
}

type GetCouponList struct {
	UserA
}

func (this *GetCouponList)GetFun() []func(){
	type Coupon struct {
		Id int64 `json:"id" `
		Name string `json:"name" `
		EndTime int64 `json:"endTime" orm:"column(endTime)"`
		Amount float64 `json:"amount" `
	}
	var(
		couponType int64
		couponList []Coupon
	)
	return []func(){
		func() {
			couponType,this.checkErr = this.ctrl.GetInt64("couponType")
			this.checkZero = couponType
		},
		func() {
			now := utils.GetTime()
			if couponType == 1{
				this.sqlQuery("select id,name,endTime,amount from coupon where isShelves = 1 and ?<endTime",&couponList,now)
			}else if couponType == 2{
				this.sqlQuery("select name,endTime,amount from userCoupon where userId=? and ?<endTime and isUsered=0",&couponList,this.user.Id,now)
			}else if couponType == 3{
				this.sqlQuery("select name,endTime,amount from userCoupon where userId=? and ? > endTime",&couponList,this.user.Id,now)
			}
		},
		func() {
			this.setReturn(1,couponList)
		},
	}
}

type GetQuestion struct {
	UserA
}

func (this *GetQuestion)GetFun() []func(){
	type Question struct {
		Question string `json:"question" `
		Answer string `json:"answer" `
	}
	var questionList []Question
	return []func(){
		func() {
			this.sqlQuery("select question,answer from question ",&questionList)
		},
		func() {
			this.setReturn(1,questionList)
		},
	}
}

type GetUserMsg struct {
	UserA
}

func (this *GetUserMsg)GetFun() []func(){
	var(
		msgNotifyList []models.ReplyMsgNotify
	)
	return []func(){
		func() {
			this.sqlQuery("select * from replyMsgNotify where srcUserId=? and isHaveRead=0",&msgNotifyList,this.user.Id)
			this.check = false;
		},
		func() {
			pwdStatus := this.user.Password != ""
			this.setReturn(1, map[string]interface{}{
				"username":this.user.Username,
				"coin":this.user.Coin,
				"img":this.user.Img,
				"pwdStatus":pwdStatus,
				"amount":this.user.Amount,
				"msgNotifyCount":len(msgNotifyList),
			})
		},
	}
}

type AvatarUpload struct {
	UserA
}

func (this *AvatarUpload)GetFun() []func(){
	var (
		url string
		cacheList []models.Cache
		userMsg map[string]map[string]interface{}
		userMsgByte []byte
	)
	return []func(){
		func() {
			url = "static/upload/user/"+utils.GetRandomString(20) +".jpg"
			this.checkErr = this.ctrl.SaveFile(url,"file")
		},
		func() {
			this.sqlQuery("select userMsg from `cache` where id = 1",&cacheList)
		},
		func() {
			 this.checkErr = json.Unmarshal([]byte(cacheList[0].UserMsg),&userMsg)
		},
		func() {
			url = "/"+url
			userMsg[this.user.Id]["img"] = url
		},
		func() {
			userMsgByte,this.checkErr = json.Marshal(userMsg)
		},
		func() {
			this.sqlTransaction(func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `user` set img=? where id=?",res,url,this.user.Id)
			}, func(o orm.Ormer, res *bool) {
				utils.SqlTranExec(o,"update `cache`  set userMsg = ? where id = 1",res,string(userMsgByte))
			})
		},
		func() {
			this.user.Img = url
			this.setReturn(1,url)
		},
	}
}

type UserSetAmount struct {
	UserA
	Amount float64
}

func (this *UserSetAmount)GetFun() []func(){
	return []func(){
		func() {
			this.checkBool = this.Amount != 0
		},
		func() {
			this.sqlExec("update `user` set amount = ? where id=?",this.Amount,this.user.Id)
		},
		func() {
			this.flushUser()
		},
	}
}

















