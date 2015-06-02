package ctrls

import (
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Dep struct {
	BaseController
}

//list of departaments
func (c *Dep) Index() {
	var deps []M.Dep
	err := DM.FindAll(&M.Dep{}, &deps, M.Sf{"Id", "Title"}, M.Where{})
	if err != nil {
		beego.Error(err)
	}
	c.Data["deps"] = deps
	c.Data["title"] = T("dep", 2)
}

//edit concrete department
func (c *Dep) Get() {
	c.Data["isEdit"] = true
	c.TplNames = "dep/form.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
		id = 0
	}
	var dep M.Dep
	if DM.FindByPk(&dep, id) != nil {
		flash := beego.NewFlash()
		flash.Notice(T("dep_not_found"))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Dep.Index"), 302)
		return
	}

	c.Data["dep"] = dep
}

//save concrete or new department
func (c *Dep) Post() {

	var v M.ValidMap
	var err error
	var msg string

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "dep/form.tpl"

	d := M.Dep{}
	flash := beego.NewFlash()
	err = c.ParseForm(&d)

	if err != nil {
		beego.Error(err)
		flash.Notice(T("internal"))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Dep.Index"), 302)
		return
	}

	if d.Id == 0 {
		v, err = DM.Insert(&d, "Title")
		msg = T("dep_created")
	} else {
		v, err = DM.Update(&d, "Title")
		msg = T("dep_updated", map[string]interface{}{"Title": d.Title})
		c.Data["isEdit"] = true
	}

	if err != nil {
		beego.Error(err)
	}

	if !v.HasErrors() {

		flash.Notice(msg)
		flash.Store(&c.Controller)

		c.Redirect(beego.UrlFor("Dep.Index"), 302)
		return
	}

	M.ExpandFormErrors(&v, c.Data)
	c.Data["dep"] = d

}

//add new department
func (c *Dep) Add() {
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["isEdit"] = false
	c.TplNames = "dep/form.tpl"
}

//remove department
func (c *Dep) Remove() {

	s := T("dep_not_found")
	d := M.Dep{}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.DeleteByPk(&d, id)
	}

	if err != nil {
		beego.Error(err)
	} else {
		s = T("dep_removed", map[string]interface{}{"Title": d.Title})
	}

	if c.IsAjax() {

		c.Data["json"] = RJson{s, err == nil}
		return

	}

	flash := beego.NewFlash()
	flash.Notice(s)
	flash.Store(&c.Controller)

	c.Redirect(beego.UrlFor("Dep.Index"), 302)
}
