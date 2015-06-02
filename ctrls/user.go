package ctrls

import (
	"database/sql"
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type User struct {
	BaseController
}

//list of users
func (c *User) Index() {

	var users []M.User

	err := DM.FindAll(&M.User{}, &users, M.Sf{}, M.Where{}, M.NewParams(M.Params{Sort: "Login ASC"}))
	if err != nil {
		beego.Error(err)
	}
	c.Data["users"] = users

	r := c.rolesList()
	m := make(map[int]string)
	for _, v := range r {
		m[v.Id] = v.Name
	}
	c.Data["roles"] = m
	c.Data["title"] = T("user", 2)
}

//view concrete user
func (c *User) Get() {

	var u M.User
	var err error

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["isEdit"] = true
	c.TplNames = "user/form.tpl"

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.FindByPk(&u, id)
	}

	if err != nil {
		flash := beego.NewFlash()
		flash.Notice(T("user_not_found"))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("User.Index"), 302)
		return
	}

	c.Data["user"] = u
	c.Data["roles"] = c.rolesList()
}

//add new user
func (c *User) Add() {
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["isEdit"] = false
	c.TplNames = "user/form.tpl"
	c.Data["roles"] = c.rolesList()
}

//add new user
func (c *User) Post() {

	var v M.ValidMap
	var err error
	var msg string

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "user/form.tpl"

	d := M.User{}

	flash := beego.NewFlash()

	err = c.ParseForm(&d)

	if err != nil {
		beego.Error(err)
		flash.Notice(T("internal"))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("User.Index"), 302)
		return
	}

	d.GToken = sql.NullString{c.Ctx.Request.Form["gtoken"][0], true}

	d.Active = true

	if d.Id == 0 {
		v, err = DM.Insert(&d)
		msg = T("user_created")
	} else {
		v, err = DM.Update(&d)
		msg = T("user_updated", map[string]interface{}{"Login": d.Login})
		c.Data["isEdit"] = true
	}

	if err != nil {
		beego.Error(err)
	}

	if !v.HasErrors() {
		flash.Notice(msg)
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("User.Index"), 302)
		return
	}

	M.ExpandFormErrors(&v, c.Data)
	c.Data["user"] = d
	c.Data["roles"] = c.rolesList()
}

//remove user
func (c *User) Remove() {

	s := T("user_not_found")
	d := M.User{}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.DeleteByPk(&d, id)
	}

	if err != nil {
		beego.Error(err)
	} else {
		s = T("user_removed", map[string]interface{}{"Login": d.Login})
	}

	if c.IsAjax() {
		c.Data["json"] = RJson{s, err == nil}
		return
	}

	flash := beego.NewFlash()
	flash.Notice(s)
	flash.Store(&c.Controller)

	c.Redirect(beego.UrlFor("User.Index"), 302)
}

func (c *User) rolesList() (r []M.Role) {
	err := DM.FindAll(&M.Role{}, &r, M.Sf{}, M.Where{})
	if err != nil {
		beego.Error(err)
	}
	return
}
