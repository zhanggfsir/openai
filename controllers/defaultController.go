package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	this.TplName = "index.tpl"
}