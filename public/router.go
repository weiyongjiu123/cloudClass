package public

import (
	"cloudClass/models"
	"cloudClass/request"
)

//不需要权限
//var TeacherRouter = map[string]string{
//	"/teacher/login":"post:Login",
//	"/teacher/testSession":"get:Session",
//	"/teacher/register":"post:Register",
//	"/teacher/getCode":"get:GetCode",
//}
//
////需要权限
//var TeacherRouterPer = map[string]string{
//	"/teacher/addClass":"POST:AddClass",
//	"/teacher/getTeacherMsg":"GET:GetTeacherMsg",
//	"/teacher/addClassInit":"GET:AddClassInit",
//	"/teacher/getClassList":"GET:GetClassList",
//	"/teacher/getChapterTitle":"GET:GetChapterTitle",
//	"/teacher/modifyChapterTitleName":"POST:ModifyChapterTitleName",
//}



var userOrdinary = map[string]interface{}{
	"/user/codeLogin":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.CodeLogin)},
	},
	"/user/getCommunication":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetCommunication)},
	},
	"/user/getHomeData":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetHomeData)},
	},
	"/user/getCatalog":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetCatalog)},
	},
	"/user/video":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.VideoPlayer)},
	},
	//获取课程评论
	"/user/getEvaluationList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetEvaluationList)},
	},
	//用户模拟登录
	"/user/UserLogin":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserLogin)},
	},
	//通过typeid获取课程列表
	"/user/UserGetClassByTypeId":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserGetClassByTypeId)},
	},
	//搜索关键字
	"/user/userClassSearch":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserClassSearch)},
	},


	//V。。。
	"/user/vGetIndex":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.VGetIndex)},
	},
	"/user/getVTypeList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetVTypeList)},
	},
	"/user/vDelTypeName":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VDelTypeName)},
	},
	"/user/vAddType":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VAddType)},
	},
	"/user/vModifyName":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VModifyName)},
	},
	"/user/vGetScopeList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.VGetScopeList)},
	},
	"/user/vAddScope":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VAddScope)},
	},
	"/user/vDelScope":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VDelScope)},
	},
	"/user/vModifyScopeName":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VModifyScopeName)},
	},
	"/user/vAddPattern":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VAddPattern)},
	},
	"/user/vDelPattern":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VDelPattern)},
	},
	"/user/vUpdatePattern":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.VUpdatePattern)},
	},
}

var userPermissions = map[string]interface{}{
	//"/user/user":map[string]interface{}{
	//	"method":"GET:test",
	//	"reqDataValid": func() interface{}{return new(request.User1)},
	//},
	"/user/user1":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.User1)},
	},
	//测试
	"/user/postTest":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.PostTest)},
	},
	"/user/getSession":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetSession)},
	},
	//添加留言
	"/user/addCommunication":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddCommunication)},
	},
	//发布留言
	"/user/commReply":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.CommReply)},
	},
	//设置用户名，更新用户表和缓存表
	"/user/setUsername":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.SetUsername)},
	},
	//回复
	"/user/commReplyOther":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.CommReplyOther)},
	},
	//获取用户个人的留言信息通知
	//"/user/getPerReplyNotify":map[string]interface{}{
	//	"method":"GET:Index",
	//	"reqDataValid": func() interface{}{return new(request.GetPerReplyNotify)},
	//},
	//测试 更新名字和用户图像路径路径
	"/user/testSetNameImg":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TestSetNameImg)},
	},
	//删除评论
	"/user/commDelete":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.CommDelete)},
	},
	//删除题主评论，在下面回复的信息也一起删除
	"/user/deleteCommunication":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DeleteCommunication)},
	},
	//获取回复通知 包括删除通知
	"/user/getReplyNotify":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetReplyNotify)},
	},
	"/user/getMsgDetail":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetMsgDetail)},
	},
	//添加学习记录
	"/user/addStudyRecord":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddStudyRecord)},
	},
	"/user/getStudyRecord":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetStudyRecord)},
	},
	"/user/addShopCar":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddShopCar)},
	},
	"/user/getShopCar":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetShopCar)},
	},
	"/user/shopCarDelete":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ShopCarDelete)},
	},
	//添加订单并支付
	"/user/addOrder":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddOrder)},
	},
	//获取课程详情页
	"/user/getClassData":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetClassData)},
	},
	"/user/readPay":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.ReadPay)},
	},
	"/user/userAddCoupon":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserAddCoupon)},
	},
	"/user/userPay":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserPay)},
	},
	"/user/getBill":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetBill)},
	},
	"/user/getOrder":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetOrder)},
	},
	"/user/addEvaluation":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddEvaluation)},
	},
	//获取一个订单，使用情景，在提交评论的页面获取
	"/user/getOneOrder":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetOneOrder)},
	},
	//签到
	"/user/sign":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.Sign)},
	},
	//获取签到记录
	"/user/getSignList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetSignList)},
	},
	//获取签到记录
	"/user/getCouponList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetCouponList)},
	},
	//获取问题列表
	"/user/getQuestion":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetQuestion)},
	},
	//获取用户信息
	"/user/getUserMsg":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetUserMsg)},
	},
	//上传用户头像
	"/user/avatarUpload":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AvatarUpload)},
	},
	//添加或更新交流留言投诉u
	"/user/userCommPlaints":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserCommPlaints)},
	},
	//支付前坚持用户已经设置支付密码
	"/user/userReadPay":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.UserReadPay)},
	},
	//更改密码
	"/user/modifyPayPwd":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ModifyPayPwd)},
	},
	//更改用户金额
	"/user/userSetAmount":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserSetAmount)},
	},
	//用户投诉
	"/user/userComplaints":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserComplaints)},
	},
}
var adminOrdinary = map[string]interface{}{
	//登录
	"/admin/adminLogin":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminLogin)},
	},
}
var adminPermissions = map[string]interface{}{
	//处理提现（已经完成提现）
	"/admin/doWithdrawal":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DoWithdrawal)},
	},
	//管理员获取主页的信息
	"/admin/adminGetIndexMsg":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetIndexMsg)},
	},
	//获取用户列表
	"/admin/adminGetUserList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetUserList)},
	},
	//禁止或解除禁止用户的登录行为
	"/admin/userBan":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UserBan)},
	},
	//获取讲师列表
	"/admin/getTeacherList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTeacherList)},
	},
	//禁止或解除讲师登录
	"/admin/teacherBan":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherBan)},
	},
	//获取投诉列表
	"/admin/getComplaintsList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetComplaintsList)},
	},
	//设置投诉为已处理或未处理
	"/admin/complaintsDo":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ComplaintsDo)},
	},
	//获取讲师提现列表
	"/admin/getWithdrawalList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetWithdrawalList)},
	},
	//设置提现失败或成功
	"/admin/withdrawalDo":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.WithdrawalDo)},
	},
	//获取课程列表
	"/admin/getClassList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetClassList)},
	},
	//禁止课程上架或运行课程上架
	"/admin/classBan":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ClassBan)},
	},
	//获取订单列表
	"/admin/adminGetOrderList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetOrderList)},
	},
	//获取一个账单
	"/admin/getBillOne":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetBillOne)},
	},
	//获取常见问题列表
	"/admin/adminGetQuestionList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetQuestionList)},
	},
	//添加问答
	"/admin/addQuestion":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddQuestion)},
	},
	//删除问答q
	"/admin/questionDel":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.QuestionDel)},
	},
	//获取类型列表
	"/admin/adminGetTypeList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetTypeList)},
	},
	//修改类型名字
	"/admin/modifyTypeName":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ModifyTypeName)},
	},
	//添加类型名字
	"/admin/addType":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddType)},
	},
	//删除类型
	"/admin/delType":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DelType)},
	},
	//获取广告列表
	"/admin/getAdList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetAdList)},
	},
	//添加广告
	"/admin/addAd":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddAd)},
	},
	//删除广告
	"/admin/delAd":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DelAd)},
	},
	//获取课程交流列表
	"/admin/adminGetQuestion":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetQuestion)},
	},
	//获取举报列表
	"/admin/getCommPlaintList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetCommPlaintList)},
	},
	//删除答疑交流（删除题主）
	"/admin/adminDeleteComm":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminDeleteComm)},
	},
	//删除答疑交流（删除一条）
	"/admin/adminDelOneComm":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminDelOneComm)},
	},
	//删除答疑交流（直接删除一条）
	"/admin/adminDelOneCommDic":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminDelOneCommDic)},
	},
	//删除答疑交流（直接删除题主）
	"/admin/adminDeleteCommDic":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminDeleteCommDic)},
	},
	//获取常规设置
	"/admin/getSetting":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetSetting)},
	},
	//修改获取常规设置
	"/admin/settingModify":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.SettingModify)},
	},
	//修改管理员的姓名图片
	"/admin/adminAddName":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminAddName)},
	},
	//获取优惠卷列表
	"/admin/adminGetCouponList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetCouponList)},
	},
	//设置优惠卷上下架状态
	"/admin/adminDoIsShelves":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminDoIsShelves)},
	},
	//添加优惠卷
	"/admin/adminAddCoupon":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AdminAddCoupon)},
	},
	//获取平台通知
	"/admin/getPlatformList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetPlatformList)},
	},
	//测试
	"/admin/test":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.Test)},
	},
	//管理员退出
	"/admin/AdminLogout":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminLogout)},
	},
	//管理员获取讲师通知
	"/admin/adminGetTeacherNotifyList":map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.AdminGetTeacherNotifyList)},
	},
	//管理员添加讲师通知
	"/admin/addTeacherNotify":map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddTeacherNotify)},
	},
}

var teacherOrdinary = map[string]interface{}{
	//讲师登录
	"/teacher/login": map[string]interface{}{
		"method":"post:Login",
	},
	//注册
	"/teacher/register": map[string]interface{}{
		"method":"post:Register",
	},
	//
	"/teacher/testSession": map[string]interface{}{
		"method":"get:Session",
	},
	//获取验证码
	"/teacher/getCode": map[string]interface{}{
		"method":"get:GetCode",
	},
	//上传视频
	//"/teacher/upload": map[string]interface{}{
	//	"method":"POST:UpdateVideo",
	//},
	//讲师登录
	//"/teacher/teacherLogin": map[string]interface{}{
	//	"method":"GET:Index",
	//	"reqDataValid": func() interface{}{return new(request.TeacherLogin)},
	//},
	//获取注册验证码
	"/teacher/teacherGetCode": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherGetCode)},
	},
	//讲师注册
	"/teacher/teacherRegister": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherRegister)},
	},
	//讲师登录
	"/teacher/TeacherLogin": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherLogin)},
	},
	//讲师登出
	"/teacher/teacherLogout": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherLogout)},
	},
	//讲师重置密码
	"/teacher/teacherResetPwd": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherResetPwd)},
	},
	//讲师重置密码发送短信
	"/teacher/teacherGetCodeResetPwd": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherGetCodeResetPwd)},
	},
}


var teacherPermissions = map[string]interface{}{
	//添加课程
	"/teacher/addClass": map[string]interface{}{
		"method":"POST:AddClass",
	},
	//获取讲师信息
	"/teacher/getTeacherMsg": map[string]interface{}{
		"method":"GET:GetTeacherMsg",
	},
	//
	"/teacher/addClassInit": map[string]interface{}{
		"method":"GET:AddClassInit",
	},
	"/teacher/getClassList": map[string]interface{}{
		"method":"GET:GetClassList",
	},
	"/teacher/getChapterTitle": map[string]interface{}{
		"method":"GET:GetChapterTitle",
	},
	"/teacher/modifyChapterTitleName": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ModifyChapterName)},
	},
	"/teacher/deleteChapter": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DeleteChapterTitle)},
	},
	"/teacher/getChapterList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetChapterList)},
	},
	"/teacher/addChapterTitle": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddChapterTitle)},
	},
	"/teacher/addChapter": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddChapter)},
	},
	"/teacher/delChapter": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DelChapter)},
	},
	"/teacher/modifyChapter": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ModifyChapter)},
	},
	"/teacher/getTypeList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTypeList)},
	},
	"/teacher/modifyClass": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.ModifyClass)},
	},
	"/teacher/delClass": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DelClass)},
	},
	"/teacher/getQuestionList": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.GetQuestionList)},
	},
	"/teacher/teacherReply": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherReply)},
	},
	"/teacher/setTeacherName": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.SetTeacherName)},
	},
	//删除评论
	"/teacher/deleteTeacherComm": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.DeleteTeacherComm)},
	},
	"/teacher/addComm": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.AddComm)},
	},
	//删除评论
	"/teacher/teacherDeleteComm": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherDeleteComm)},
	},
	"/teacher/getNotifyList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetNotifyList)},
	},
	"/teacher/getNotifyDetail": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetNotifyDetail)},
	},
	"/teacher/getOrderList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetOrderList)},
	},
	"/teacher/getTeacherMessage": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTeacherMessage)},
	},
	"/teacher/setTeacherPhone": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.SetTeacherPhone)},
	},
	//提现
	"/teacher/withdrawal": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.Withdrawal)},
	},
	//获取讲师的账单
	"/teacher/getTeacherBill": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTeacherBill)},
	},
	//获取讲师首页的数据与信息
	"/teacher/getTeacherIndexMsg": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTeacherIndexMsg)},
	},
	//获取讲师的通知（包括平台发的通知）
	"/teacher/getTeacherNotifyList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.GetTeacherNotifyList)},
	},
	//投诉留言
	"/teacher/teacherCommPlaints": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherCommPlaints)},
	},
	//获取讲师通知列表
	"/teacher/teacherGetNotifyList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherGetNotifyList)},
	},
	//获取讲师通知列表(一条)
	"/teacher/getOneTeacherNotify": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.GetOneTeacherNotify)},
	},
	//讲师获取课程列表
	"/teacher/teacherGetClassList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherGetClassList)},
	},
	//获取课程类型列表
	"/teacher/teacherGetTypeList": map[string]interface{}{
		"method":"GET:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherGetTypeList)},
	},
	//更改课程上下架状态
	"/teacher/updateIsShelves": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.UpdateIsShelves)},
	},
	//删除讲师通知
	"/teacher/teacherNotifyDel": map[string]interface{}{
		"method":"POST:Index",
		"reqDataValid": func() interface{}{return new(request.TeacherNotifyDel)},
	},
}

var RouterMap = map[string]map[string]interface{}{
	"teacher":{
		"ordinary":teacherOrdinary,
		"permissions":teacherPermissions,
	},
	"user":{
		"ordinary":userOrdinary,
		"permissions":userPermissions,
	},
	"admin":{
		"permissions":adminPermissions,
		"ordinary":adminOrdinary,
	},
}

var ModuleSessionMap = map[string]map[string]interface{}{
	"teacher":{
		"name":"teacher",
		"returnSession": func(name string,session interface{}) map[string]interface{}{
			teacher := session.(models.Teacher)
			return map[string]interface{}{
				"amount":teacher.Amount,
				"img":teacher.Img,
				"name":teacher.Name,
				"phone":teacher.Phone,
				"createTime":teacher.CreateTime,
				"id":teacher.Id,
			}
		},
	},
	"user":{
		"name":"user",
		"returnSession": func(name string,session interface{}) map[string]interface{}{ //小程序不行需要返回用户信息
			return map[string]interface{}{}
		},
	},
	"admin":{
		"name":"admin",
		"returnSession": func(name string,session interface{}) map[string]interface{}{ //小程序不行需要返回用户信息
		return map[string]interface{}{}
	},
	},
}
