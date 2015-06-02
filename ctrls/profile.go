package ctrls

import (
	"html/template"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Prof struct {
	BaseController
}

//show edit account settings form
func (c *Prof) Get() {
	var u M.User
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())

	err := DM.FindByPk(&u, c.GetSession("uid").(int))
	if err != nil {
		beego.Error(err)
		return
	}

	c.Data["Name"] = u.Name
	c.Data["user"] = u
}

//save account settings
func (c *Prof) Post() {
	var v M.ValidMap
	var err error
	var msg string
	s := []string{"Name", "Email", "NotifyByMail", "NotifyByTelegram"}

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "prof/get.tpl"

	d := M.Profile{}
	err = c.ParseForm(&d)
	d.Id = c.GetSession("uid").(int)

	if err == nil && d.Password != "" {
		s = append(s, "Password")
	}

	if err == nil {
		v, err = DM.Update(&d, s...)
	}

	if err == nil && !v.HasErrors() {
		c.SetSession("name", d.Name)
		msg = T("account_saved")
	}

	if err != nil {
		msg = T("internal")
		beego.Error(err)
	}

	if msg != "" {
		flash := beego.NewFlash()
		flash.Notice(msg)
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Prof.Get"), 302)
		return
	}

	M.ExpandFormErrors(&v, c.Data)
	c.Data["Name"] = d.Name
}
