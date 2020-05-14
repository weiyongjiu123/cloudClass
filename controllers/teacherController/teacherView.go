package teacherController

func (this *TeacherController) IndexView() {
	this.TplName = "teacher/index.html"
}

func (this *TeacherController)ClassManagerView()  {
	this.TplName = "teacher/classManager.html"
}

func (this *TeacherController) ChapterView(){
	this.TplName = "teacher/chapter.html"
}

func (this *TeacherController) ChapterAddView() {
	this.TplName = "teacher/chapterAdd.html"
}

func (this *TeacherController) TestView() {
	this.TplName = "teacher/test.html"
}

func (this *TeacherController) AddClassView() {
	this.TplName = "teacher/addClass.html"
}

func (this *TeacherController) ChapterTitleView() {
	this.TplName = "teacher/chapterTitle.html"
}
