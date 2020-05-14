package request

type Ctrl struct {
	GetSessionFun func(name string)interface{}
	SetSessionFun func(name interface{},content interface{})
	GetParStringFun func(name string)string
	GetParIntFun func(name string)int
	GetInt64 func(name string,def ...int64) (int64,error)
	SaveFile func(url string,key string)error
	Redirect func(url string,code int)
}

func (this *Ctrl)SetSession(name interface{},content interface{}){
	this.SetSessionFun(name,content)
}

func (this *Ctrl)GetSession(name string)interface{}{
	return this.GetSessionFun(name)
}

func (this *Ctrl) GetString(name string)string {
	return this.GetParStringFun(name)
}

func (this *Ctrl) GetInt(name string)int {
	return this.GetParIntFun(name)
}

//func (this *Ctrl) GetInt64(name string)int64 {
//	return this.GetParInt64Fun(name)
//}


