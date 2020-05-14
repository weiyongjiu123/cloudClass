package request

type GetCommunicationRes struct {
	UserType int8 `json:"userType" orm:"column(userType)"`
	UserId string `json:"userId" orm:"column(userId)"`
	Username string `json:"username" orm:"column(username)"`
	Img string `json:"img" orm:"column(img)"`
	Content string `json:"content" orm:"column(content)"`
	CreateTime int64 `json:"createTime" orm:"column(createTime)"`
	ReplyContent string `json:"-" orm:"column(replyContent)"`
	ClassId int64 `json:"classId" orm:"column(classId)"`
	CommId int64 `json:"commId" orm:"column(id)"`
	ReplyContentStrut map[string]ReplyContent `json:"replyList" `
	ReplyContentHeader string `json:"replyContentHeader" orm:"column(replyContentHeader)"`
	ReplyContentTail string `json:"replyContentTail" orm:"column(replyContentTail)"`
}


type ReplyContent struct {
	PostUserId string `json:"postUserId"`
	SrcUserId string `json:"srcUserId"`
	Content string `json:"content"`
	CreateTime int64 `json:"createTime"`
	PostUsername string `json:"postUsername"`
	SrcUsername string `json:"srcUsername"`
	Next string `json:"next"`
	Pre string `json:"pre"`
	PostUserType int8 `json:"postUserType"`
	SrcUserType int8 `json:"srcUserType"`
}

type Cache struct {
	VideoPer string  `json:"videoPer" orm:"column(videoPer)"`
	UserMsg string `json:"userMsg" orm:"column(userMsg)"`
}

type VideoPer struct {
	Time int64
	Url string
}





