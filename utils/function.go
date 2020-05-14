package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str) )
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetTime() int64 {
	return time.Now().Unix()
}

func GetDate(t int64) string{
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	datetime := time.Unix(t, 0).Format(timeLayout)
	return datetime
}

func GetNowDate()string  {
	now := time.Now()
	year := now.Year() //年
	month :=now.Month()
	day := now.Day()
	return strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
}

func GetNowDateA()string  {
	var dayStr string
	now := time.Now()
	year := now.Year() //年
	month :=now.Month()
	day := now.Day()
	if day < 10{
		dayStr = "0"+strconv.Itoa(day)
	}else{
		dayStr = strconv.Itoa(day)
	}
	return strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + dayStr
}


func GetYesterdayData() string{
	nTime := time.Now()
	yesTime := nTime.AddDate(0,0,-1)
	return yesTime.Format("2006-1-02")
}

func GetYesterdayDataA() string{
	nTime := time.Now()
	yesTime := nTime.AddDate(0,0,-1)
	return yesTime.Format("2006-1-2")
}



func GetDateBy(t int64,format string)string {
	//timeLayout := "2006-01-02"  //转化所需模板
	return time.Unix(t, 0).Format(format)
}

var randNumber int64
//随机数
func RandNum(min int64, max int64)  int64{
	if randNumber == 0{
		rand.Seed(time.Now().UnixNano())
	}else{
		rand.Seed(randNumber)
	}

	num := rand.Int63n(max-min)
	num += min
	randNumber = num
	return num
}

//随机字符串
func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func SendCode(phone string) (string,bool) {
	fmt.Println("send code to phone ...")
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(8999)
	code += 1000
	fmt.Println("generate code",code)
	return strconv.Itoa(code),true
}

func SqlQuery(sql string, params ...interface{})  ([]orm.Params,int64){
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	var maps []orm.Params
	o :=orm.NewOrm()
	num ,error := o.Raw(sql,params).Values(&maps)
	if error != nil{
		fmt.Println(error)
	}
	return maps,num
}

func SqlTransaction(fun...func(o orm.Ormer,res *bool)) bool {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil{
		if sqlPrint{
			fmt.Println("sql:",err)
		}
		return false
	}
	for _, v := range fun {
		res :=true;
		v(o,&res)
		if !res{
			o.Rollback();
			return false
		}
	}
	o.Commit()
	return true
}



func SqlTranExec(ormer orm.Ormer, sql string, resBool *bool,params ...interface{})(lastId int64) {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	res,err := ormer.Raw(sql,params).Exec()

	if err != nil{
		fmt.Println(err.Error())
		*resBool = false
		return
	}else{
		affected,_ := res.RowsAffected()
		lastId,_ = res.LastInsertId()
		if affected > 0{
			*resBool = true
		}else{
			*resBool = false
		}
	}
	return
}

func SqlExec(sql string,params ... interface{})(isSuccess bool,lastId int64,affected int64) {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	o := orm.NewOrm()
	res,err := o.Raw(sql,params).Exec()
	if err != nil{
		isSuccess = false
		fmt.Println(err.Error())
		return
	}
	lastId,_ = res.LastInsertId()
	affected,_ = res.RowsAffected()
	if affected > 0{
		isSuccess = true
	}else{
		isSuccess = false
	}
	return
}

//合并map
func MergeMap(mapDst* map[string]string ,mapSrc map[string]string)  {
	for k, v := range mapSrc {
		(*mapDst)[k] = v
	}
}

//字符串连接
func StrMerge(s1 string,s2 string) string {
	var stringBuilder bytes.Buffer
	stringBuilder.WriteString(s1)
	stringBuilder.WriteString(s2)
	return stringBuilder.String()
}

//字符串截取
func StrSeg(srcStr string, s string)  (arr []string,b bool){
	re := regexp.MustCompile(s)
	res := re.FindStringSubmatch(srcStr)
	if len(res) > 0{
		b = true
		arr = res
	}else{
		b = false
	}
	return
}

func SaveBase64ToPic(picBase64 string,picUrl string)  bool{
	re := regexp.MustCompile("base64,([\\s\\S]*)$")
	res := re.FindStringSubmatch(picBase64)
	if len(res) == 0 {
		return false
	}
	ddd, _ := base64.StdEncoding.DecodeString(res[1])
	err2 := ioutil.WriteFile(picUrl, ddd, 0666)
	if err2 != nil{
		return false
	}
	return true
}

func TestMap(list []orm.Params) {
	for _, v := range list {
		fmt.Println(v)
	}
}

func getDataType(s string,o interface{})  {
	switch s {
	case "int":

	}
}


func Typeof(v interface{})string {
	return reflect.TypeOf(v).String()
}

func GetChapterTitleNumStr(num int)  string{
	arr := []string{"第一章", "第二章", "第三章", "第四章", "第五章", "第六章", "第七章", "第八章", "第九章","第十章", "第十一章", "第十二章", "第十三章", "第十四章", "第十五章", "第十六章", "第十七章", "第十八章", "第十九章", "第二十章"}
	return arr[num]
}

func SavePicByBase64(base64 string) (bool,string) {
	t := GetTime()
	randNum := RandNum(100000, 999999)
	picName := StrMerge(strconv.FormatInt(t,10),strconv.FormatInt(randNum,10))
	picUrl := StrMerge(beego.AppConfig.String("classimg"),picName)
	picUrl = StrMerge(picUrl,".jpg")
	return SaveBase64ToPic(base64,picUrl),"/"+picUrl
}

//获取微信小程序用户openId
func GetOpenId(code string) (ok bool,openId string){
	appId := beego.AppConfig.String("wxappid")
	appSecret := beego.AppConfig.String("wxsecret")
	url := "https://api.weixin.qq.com/sns/jscode2session?appid="+appId+"&secret="+appSecret+"&js_code="+code+"&grant_type=authorization_code";
	req := httplib.Get(url)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	str,err := req.String()
	if err != nil{
		fmt.Println(err)
		ok = false
		return
	}
	data := make(map[string]interface{})
	err =  json.Unmarshal([]byte(str),&data)
	if err != nil{
		fmt.Println(err)
		ok = false
		return
	}
	ok = true
	openId = data["openid"].(string)
	return
}

func SqlQueryStrut(sql string,arr interface{},params...interface{}) (int64,error)  {
	sqlPrint,_ := beego.AppConfig.Bool("sqlprint")
	if sqlPrint{
		fmt.Println("sql:",sql,params)
	}
	o := orm.NewOrm()
	return  o.Raw(sql,params).QueryRows(arr)
}

func ValIsExistArr(arr []string ,value string) bool {
	for e := range arr {
		if arr[e] == value {
			return true
		}
	}
	return false
}

type Page struct {
	Text int `json:"text"`
	Start int  `json:"start"`
	Now int8 `json:"now"`
}

func GetPage(count int,start int) []Page{
	var(
		pageList []Page
		number int
		numCount int
	)
	number = start / 10 +1
	if count % 10 == 0{
		numCount = count / 10
	}else{
		numCount = count / 10 + 1
	}
	if number  <= 6{
		number = 1
	}else if number + 6 >= numCount{
		number =  numCount - 9
	}else{
		number = number -4
	}
	num:= 1
	for i := (number-1) * 10;i < count ;i+=10  {
		var now int8
		if start == i{
			now = 1
		}else{
			now = 0
		}
		pageList = append(pageList,Page{number,i,now})
		if num == 10{
			break
		}
		number++
		num++
	}
	return  pageList
}

func SubStrRunes(s string, length int) string {
	if utf8.RuneCountInString(s) > length {
		rs := []rune(s)
		return string(rs[:length])
	}

	return s
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}