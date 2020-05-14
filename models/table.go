package models

type Class struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	Detail string `json:"detail" orm:"column(detail)"`
	IsCharge int8 `json:"isCharge" orm:"column(isCharge)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	TeacherId int64 `json:"teacherId" orm:"column(teacherId)"`
	TypeId int64 `json:"typeId" orm:"column(typeId)"`
	ClassImg string `json:"classImg" orm:"column(classImg)"`
	WatchCount int64 `json:"watchCount" orm:"column(watchCount)"`
	IsShelves int8 `json:"isShelves" orm:"column(isShelves)"`
	ChapterCount int8 `json:"chapterCount" orm:"column(chapterCount)"`
	ChapterTitleCount int8 `json:"chapterTitleCount" orm:"column(chapterTitleCount)"`
	Price float64 `json:"price" orm:"column(price)"`
	Score float32  `json:"score" orm:"column(score)"`
	Star float32  `json:"star" orm:"column(star)"`
	IsBan int8 `json:"isBan" orm:"column(isBan)"`
	UserWatchRecord string `json:"userWatchRecord" orm:"column(userWatchRecord)"`
}

type Teacher struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	Phone string `json:"phone" orm:"column(phone)"`
	Password string `json:"-" orm:"column(password)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	Img string `json:"img" orm:"column(img)"`
	UnMoveAmount float64 `json:"unMoveAmount" orm:"column(unMoveAmount)"`
	SellCount int64 `json:"sellCount" orm:"column(sellCount)"`
	WatchCount int64  `json:"watchCount" orm:"column(watchCount)"`
	IsBan int8   `json:"isBan" orm:"column(isBan)"`
}

type Chapter struct {
	Id int64 		`json:"id" orm:"column(id)"`
	ClassId int64	 `json:"classId" orm:"column(classId)"`
	Name string 	`json:"name" orm:"column(name)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	WatchCount int64 `json:"watchCount" orm:"column(watchCount)"`
	ChapterTitleId int64 `json:"chapterTitleId" orm:"column(chapterTitleId)"`
	TeacherId int64 `json:"teacherId" orm:"column(teacherId)"`
	Time string `json:"time" orm:"column(time)"`
	Video string `json:"video" orm:"column(video)"`
}

type ChapterTitle struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	ChapterCount int64 `json:"chapterCount" orm:"column(chapterCount)"`
	TeacherId int64 `json:"teacherId" orm:"column(teacherId)"`
}

type Type struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	ModifyTime int64 `json:"modifyTime" orm:"column(modifyTime)"`
}

type User struct {
	Id string `json:"id" orm:"column(id)"`
	Img string `json:"img" orm:"column(img)"`
	Coin float64 `json:"coin" orm:"column(coin)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	Username string `json:"username" orm:"column(username)"`
	Phone string `json:"phone" orm:"column(phone)"`
	Amount float64   `json:"amount"`
	Password string `json:"password" orm:"column(password)"`
	IsBan int8  `json:"isBan" orm:"column(isBan)"`
}

type Communication struct {
	Id int64 `json:"id" orm:"column(id)"`
	UserId string `json:"userId" orm:"column(userId)"`
	Content string `json:"content" orm:"column(content)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	ReplyContent string `json:"replyContent" orm:"column(replyContent)"`
	ClassId int64  `json:"classId" orm:"column(classId)"`
	ReplyContentHeader string `json:"replyContentHeader" orm:"column(replyContentHeader)"`
	ReplyContentTail string `json:"replyContentTail" orm:"column(replyContentTail)"`
	UserType int8 `json:"userType" orm:"column(userType)"`
}

type Cache struct {
	Id int64 `json:"id" orm:"column(id)"`
	UserMsg string `json:"userMsg" orm:"column(userMsg)"`
}

type ReplyMsgNotify struct {
	Id int64 `json:"id" orm:"column(id)"`
	SrcUserId string `json:"srcUserId" orm:"column(srcUserId)"`
	SrcContent string `json:"srcContent" orm:"column(srcContent)"`
	SrcCreateTime int64 `json:"srcCreateTime" orm:"column(srcCreateTime)"`
	PostUserId string `json:"postUserId" orm:"column(postUserId)"`
	PostCreateTime int64 `json:"postCreateTime" orm:"column(postCreateTime)"`
	PostContent string `json:"postContent" orm:"column(postContent)"`
	IsHaveRead int8 `json:"isHaveRead" orm:"column(isHaveRead)"`
	NotifyType int8 `json:"notifyType" orm:"column(notifyType)"`
	ReplyIndex string  `json:"replyIndex" orm:"column(replyIndex)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	CommId int64 `json:"commId" orm:"column(commId)"`
	SrcUserType int8  `json:"srcUserType" orm:"column(srcUserType)"`
	PostUserType int8 `json:"postUserType" orm:"column(postUserType)"`

	SrcUsername string `json:"srcUsername" `
	SrcUserImg string `json:"srcUserImg" `
	PostUsername string `json:"postUsername" `
	PostUserImg string `json:"postImg" `
}

type ShopCar struct {
	Id int64  `json:"id" orm:"column(id)"`
	UserId string `json:"userId" orm:"column(userId)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}

type Order struct {
	OrderId string `json:"orderId" orm:"column(orderId)"`
	IsPay int8 `json:"isPay" orm:"column(isPay)"`
	Status int8 `json:"status" orm:"column(status)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	UserId string `json:"userId" orm:"column(userId)"`
	ClassName string `json:"className" orm:"column(className)"`
	ClassImg string `json:"classImg" orm:"column(classImg)"`
	TeacherId float64	`json:"teacherId" orm:"column(teacherId)"`
	BillId string `json:"billId" orm:"column(billId)"`
}

type UserCoupon struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	StartTime int64 `json:"startTime" orm:"column(startTime)"`
	EndTime int64 `json:"endTime" orm:"column(endTime)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	UserId string `json:"userId" orm:"column(userId)"`
	IsUsered int8 `json:"isUsered" orm:"column(isUsered)"`
}

type Coupon struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	StartTime int64 `json:"startTime" orm:"column(startTime)"`
	EndTime int64  `json:"endTime" orm:"column(endTime)"`
	IsShelves int8 `json:"isShelves" orm:"column(isShelves)"`
	ModifyAble int8  `json:"modifyAble" orm:"column(modifyAble)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	Count int64  `json:"count" orm:"column(count)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}

type Setting struct {
	Id int64 `json:"id" orm:"column(id)"`
	DealRatio float64 `json:"dealRatio" orm:"column(dealRatio)"`
	CoinRatio float64 `json:"coinRatio" orm:"column(coinRatio)"`
	SignCoin  int64 `json:"signCoin" orm:"column(signCoin)"`
}

type Bill struct {
	Id string  `json:"id" orm:"column(id)"`
	UserId string `json:"userId" orm:"column(userId)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	CouponList string `json:"couponList" orm:"column(couponList)"`
	ClassList string `json:"classList" orm:"column(classList)"`
	Coin float64 `json:"coin" orm:"column(coin)"`
	Pay float64 `json:"pay" orm:"column(pay)"`
	CoinRatio float64  `json:"coinRatio" orm:"column(coinRatio)"`
	OrderIdList string `json:"orderIdList" orm:"column(orderIdList)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
}

type Sign struct {
	UserId string  `json:"userId" orm:"column(userId)"`
	Date string `json:"date" orm:"column(date)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
}

type Admin struct {
	Id int8  `json:"id" orm:"column(id)"`
	Name string  `json:"name" orm:"column(name)"`
	Password string   `json:"password" orm:"column(password)"`
}

type Withdrawal struct {
	Id int64  `json:"id" orm:"column(id)"`
	Username string  `json:"username" orm:"column(username)"`
	Number string  `json:"number" orm:"column(number)"`
	TeacherId int64  `json:"teacherId" orm:"column(teacherId)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
	DoTime int64  `json:"doTime" orm:"column(doTime)"`
	IsDo int8  `json:"isDo" orm:"column(isDo)"`
	Amount float64	`json:"amount" orm:"column(amount)"`
	Bank string  `json:"bank" orm:"column(bank)"`
}

type TeacherBill struct {
	Id int64  `json:"id" orm:"column(id)"`
	Description string `json:"description" orm:"column(description)"`
	Amount float64 `json:"amount" orm:"column(amount)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	OrderId string `json:"orderId" orm:"column(orderId)"`
}

type WatchRecord struct {
	Id int64  `json:"id" orm:"column(id)"`
	Date string `json:"date" orm:"column(date)"`
	WatchCount int64 `json:"watchCount" orm:"column(watchCount)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	TeacherId int64 `json:"TeacherId" orm:"column(TeacherId)"`
}

type TeacherNotify struct {
	Id int64 `json:"id" orm:"column(id)"`
	Title string `json:"title" orm:"column(title)"`
	Content string `json:"content" orm:"column(content)"`
	TeacherId int64 `json:"teacherId" orm:"column(teacherId)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}

type Complaints struct {
	Id int64 `json:"id" orm:"column(id)"`
	Content string `json:"content" orm:"column(content)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	UserId string `json:"userId" orm:"column(userId)"`
	IsDo int8 `json:"isDo" orm:"column(isDo)"`
	ContentPre string `json:"contentPre" orm:"column(contentPre)"`
}

type Question struct {
	Id int64 `json:"id" orm:"column(id)"`
	Question string `json:"question" orm:"column(question)"`
	Answer string `json:"answer" orm:"column(answer)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
}
type Ad struct {
	Id int64 `json:"id" orm:"column(id)"`
	PicUrl string `json:"picUrl" orm:"column(picUrl)"`
	ClassId int64  `json:"classId" orm:"column(classId)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
}

type CommComplaints struct {
	Id string `json:"id" orm:"column(id)"`
	CommId int64 `json:"commId" orm:"column(commId)"`
	CommIndex string `json:"commIndex" orm:"column(commIndex)"`
	Content string `json:"content" orm:"column(content)"`
	IsDel int8 `json:"isDel" orm:"column(isDel)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	Count int64 `json:"count" orm:"column(count)"`
	UserId string `json:"userId" orm:"column(userId)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	UserType int8 `json:"userType" orm:"column(userType)"`
	AllContent string  `json:"allContent" orm:"column(allContent)"`
}

type PlatformNotify struct {
	Id int64 `json:"id" orm:"column(id)"`
	Title string `json:"title" orm:"column(title)"`
	Content string `json:"content" orm:"column(content)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	Type string `json:"type" orm:"column(type)"`
	LinkId int64  `json:"linkId" orm:"column(linkId)"`
}

type VPattern struct {
	Id int64 `json:"id" orm:"column(id)"`
	Pattern string `json:"pattern" orm:"column(pattern)"`
	Example string `json:"example" orm:"column(example)"`
	ScopeId int64 `json:"scopeId" orm:"column(scopeId)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}

type VScope struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	TypeId int64  `json:"typeId" orm:"column(typeId)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}

type VType struct {
	Id int64 `json:"id" orm:"column(id)"`
	Name string  `json:"name" orm:"column(name)"`
	CreateTime int64  `json:"createTime" orm:"column(createTime)"`
}
