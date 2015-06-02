package ctrls

import (
	"html/template"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Main struct {
	BaseController
}

type LoginForm struct {
	Login    string `valid:"Required"form:"login"`
	Password string `valid:"Required"form:"password"`
}

//show auth form
func (c *Main) Get() {
	if c.GetSession("uid") != nil {
		c.Redirect(beego.UrlFor("User.Index"), 302)
	}
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "loginform.tpl"
	c.Data["title"] = "HRKB"
}

//check auth form
func (c *Main) Post() {

	var v M.ValidMap

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "loginform.tpl"

	loginform := LoginForm{}
	c.ParseForm(&loginform)
	v.Valid(loginform)

	c.Data["login"] = loginform.Login
	M.ExpandFormErrors(&v, c.Data)

	if v.HasErrors() {
		return
	}

	c.Data["errLogin"] = T("login_error")
	id, err := M.CheckPass(loginform.Login, loginform.Password)

	if err != nil {
		beego.Warn(err)
		return
	}

	u := M.User{}
	err = DM.FindByPk(&u, id)
	if err != nil {
		beego.Error(err)
		return
	}

	if len(u.GToken.String) > 0 {
		c.SetSession("gitlabToken", u.GToken.String)
	}
	c.SetSession("uid", id)
	c.SetSession("role", u.Role)
	c.SetSession("name", u.Name)
	c.Redirect(beego.UrlFor("Cand.Index"), 302)
}

func (c *Main) Logout() {
	c.DestroySession()
	c.Redirect(beego.UrlFor("Main.Get"), 302)
}
