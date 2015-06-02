package ctrls

import (
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Crit struct {
	BaseController
}

// show all criteria
func (c *Crit) Index() {
	crit := &M.Crit{}
	groups, err := crit.GetGroupedCrits()
	if err != nil {
		beego.Warn("GroupError", err)
	}
	c.Data["groups"] = groups
}

// save criteria
func (c *Crit) Post() {
	var v M.ValidMap
	var err error
	msg := T("crit_created")

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "crit/form.tpl"

	flash := beego.NewFlash()
	crit := M.Crit{}

	err = c.ParseForm(&crit)

	if crit.Id != 0 {
		msg = T("crit_updated", map[string]interface{}{"Title": crit.Title})
		c.Data["isEdit"] = true
	}

	if err == nil {
		if crit.Id == 0 {
			v, err = DM.Insert(&crit, "Title", "Dep")
		} else {
			v, err = DM.Update(&crit, "Title", "Dep")
		}
	}

	if err != nil {
		beego.Error(err)
	}

	if !v.HasErrors() && err == nil {
		flash.Notice(msg)
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Crit.Index"), 302)
		return
	}

	var deps []M.Dep
	err = DM.FindAll(&M.Dep{}, &deps, M.Sf{"Id", "Title"}, M.Where{})
	if err != nil {
		beego.Error(err)
	}

	c.Data["deps"] = deps
	c.Data["crit"] = crit
	M.ExpandFormErrors(&v, c.Data)

}

// update criteria
func (c *Crit) Get() {
	c.Data["isEdit"] = true
	c.TplNames = "crit/form.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
		id = 0
	}
	var crit M.Crit
	if DM.FindByPk(&crit, id) != nil {
		flash := beego.NewFlash()
		flash.Notice(T("crit_not_found"))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Crit.Index"), 302)
		return
	}
	var deps []M.Dep
	err = DM.FindAll(&M.Dep{}, &deps, M.Sf{"Id", "Title"}, M.Where{})
	if err != nil {
		beego.Error(err)
	}

	c.Data["deps"] = deps
	c.Data["crit"] = crit
}

//add new criteria
func (c *Crit) Add() {
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.TplNames = "crit/form.tpl"
	var deps []M.Dep
	err := DM.FindAll(&M.Dep{}, &deps, M.Sf{"Id", "Title"}, M.Where{})
	if err != nil {
		beego.Error(err)
	}

	c.Data["deps"] = deps
}

//remove criteria
func (c *Crit) Remove() {
	s := T("crit_not_found")
	crit := M.Crit{}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.DeleteByPk(&crit, id)
	}

	if err != nil {
		beego.Error(err)
	} else {
		s = T("crit_removed", map[string]interface{}{"Title": crit.Title})
	}

	if c.IsAjax() {
		c.Data["json"] = RJson{s, err == nil}
		return
	}

	flash := beego.NewFlash()
	flash.Notice(s)
	flash.Store(&c.Controller)

	c.Redirect(beego.UrlFor("Crit.Index"), 302)

}
