package ctrls

import (
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Contact struct {
	BaseController
}

//list of candidates
func (c *Contact) Index() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
	}

	var contacts []M.Contact
	var cand M.Cand

	if err := DM.FindAll(&M.Contact{}, &contacts, M.Sf{}, M.Where{And: M.W{"Cand": id}}); err != nil {
		beego.Error(err)
	}

	if err := DM.Find(&cand, M.Sf{}, M.Where{And: M.W{"Id": id}}); err != nil {
		beego.Error(err)
	}

	c.Data["contacts"] = contacts
	c.Data["cand"] = cand
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
}

//add new candidate
func (c *Contact) Add() {
	var err error

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
		return
	}

	cont := M.Contact{Cand: id, Active: true}

	err = c.ParseForm(&cont)

	if err != nil {
		c.Data["json"] = struct {
			Error string `json:"error"`
		}{T("parse_error")}
		beego.Error(err)
		return
	}

	if _, err := DM.Insert(&cont); err != nil {
		c.Data["json"] = struct {
			Error string `json:"error"`
		}{T("parse_error")}
		beego.Error(err)
		return
	}

	c.Data["json"] = struct {
		Id int `json:"id"`
	}{cont.Id}
}

//remove department
func (c *Contact) Remove() {

	var s string

	id, err := strconv.Atoi(c.Ctx.Input.Param(":cid"))
	if err != nil {
		beego.Error(err)
		s = T("invalid_param", map[string]interface{}{"Param": "ID"})
	}

	d := M.Contact{}

	err = DM.DeleteByPk(&d, id)

	if err != nil {
		s = T("contact_not_found")
		beego.Error(err)
	} else {
		s = T("contact_created")
	}

	c.Data["json"] = RJson{s, err == nil}
}
