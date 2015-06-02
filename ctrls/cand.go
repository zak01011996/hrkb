package ctrls

import (
	"database/sql"
	"fmt"
	"html/template"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	M "hrkb/models"
	"hrkb/utils"
)

type Cand struct {
	BaseController
}

//list of candidates
func (c *Cand) Index() {
	var cands []M.Cand
	var deps []M.Dep
	var wF M.Where

	depF, err := strconv.Atoi(c.Input().Get("dep"))

	if depF > 0 {
		wF = M.Where{And: M.W{"Dep": depF}}
	}

	cnt, err := DM.Count(&M.Cand{}, wF)
	if err != nil {
		beego.Error(err)
	}

	perPage := 10
	paginator := pagination.SetPaginator(c.Ctx, perPage, cnt)

	//TODO properly handle db errors
	err = DM.FindAll(&M.Cand{}, &cands, M.Sf{}, wF, M.NewParams(M.Params{Offset: paginator.Offset(), Limit: perPage, Sort: "Name,LName ASC"}))
	if err != nil {
		beego.Error(err)
		return
	}
	err = DM.FindAll(&M.Dep{}, &deps, M.Sf{}, M.Where{}, M.NewParams(M.Params{Sort: "Title ASC"}))
	if err != nil {
		beego.Error(err)
		return
	}
	m := make(map[int]string)
	for _, v := range deps {
		m[v.Id] = v.Title
	}
	c.Data["deps"] = m
	c.Data["depFilter"] = depF
	c.Data["cands"] = cands
	c.Data["title"] = T("cand", 2)
}

//get concrete candidate
func (c *Cand) Get() {

	var comments []M.Comment
	var tComments []tplComment
	var contacts []M.Contact
	var cand M.Cand
	var files []M.File
	var deps []M.Dep

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())

	if err != nil {
		beego.Error(err)
	}

	if DM.FindByPk(&cand, id) != nil {
		flash := beego.NewFlash()
		flash.Notice(T("nocand", 1))
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Cand.Index"), 302)
		return
	}

	err = DM.FindAll(&M.Comment{}, &comments, M.Sf{}, M.Where{And: M.W{"Cand": id}})
	if err != nil {
		beego.Error(err)
	}

	if err = DM.FindAll(&M.Contact{}, &contacts, M.Sf{}, M.Where{And: M.W{"Cand": id}}); err != nil {
		beego.Error(err)
	}

	if err := DM.FindAll(&M.File{}, &files, M.Sf{}, M.Where{And: M.W{"Cand": id}}); err != nil {
		beego.Error(err)
	}

	err = DM.FindAll(&M.Dep{}, &deps, M.Sf{}, M.Where{})
	if err != nil {
		beego.Error(err)
	}

	for _, v := range comments {
		cr := tplComment{Id: v.Id, User: v.User, Comment: v.Comment}
		cr.SetAuthor(v.User)
		cr.SetDate(v.CreatedAt)
		tComments = append(tComments, cr)
	}

	m := make(map[int]string)
	for _, v := range deps {
		m[v.Id] = v.Title
	}

	c.Data["deps"] = m
	c.Data["Comments"] = tComments
	c.Data["contacts"] = contacts
	c.Data["files"] = files

	cr := &M.Crit{}
	crits, err := cr.GetGroupedCrits()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["crits"] = crits
	}
	c.Data["cand"] = cand
}

//save concrete or new candidate
func (c *Cand) Post() {
	upload_dir := beego.AppConfig.String("upload_dir")

	c.TplNames = "cand/add.tpl"
	flash := beego.NewFlash()

	var v M.ValidMap
	var err error
	d := M.Cand{}
	var deps []M.Dep

	if err := c.ParseForm(&d); err != nil {
		beego.Error(err)

		flash.Error(T("internal"))
		flash.Store(&c.Controller)
		return
	}

	d.Note = sql.NullString{c.Ctx.Request.Form["note"][0], true}
	d.Active = true

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["cand"] = d

	if err := DM.FindAll(&M.Dep{}, &deps, M.Sf{}, M.Where{}); err != nil {
		beego.Error("Department find error: ", err)

		flash.Error(T("internal"))
		flash.Store(&c.Controller)

		return
	}

	c.Data["deps"] = deps

	var header *multipart.FileHeader

	_, header, err = c.GetFile("img")

	if err := os.MkdirAll(upload_dir+"img/", os.ModePerm); err != nil {
		beego.Error(err)
	}

	if err := os.MkdirAll(beego.AppConfig.String("tmp_dir"), os.ModePerm); err != nil {
		beego.Error(err)
	}

	d.Img = sql.NullString{beego.AppConfig.String("static_dir") + "img/noavatar-" + strconv.Itoa(utils.RandInt(0, 4)) + ".jpg", true}

	if err == nil {

		if !strings.Contains(header.Header["Content-Type"][0], "image") {
			flash.Error(T("filetype_bad"))
			flash.Store(&c.Controller)
			return
		}

		var filename = beego.AppConfig.String("tmp_dir") + header.Filename
		c.SaveToFile("img", filename)

		furl, fsname := utils.Resize(filename, header, 350, 0)

		if furl == "" {
			flash.Error(T("resize_bad"))
			flash.Store(&c.Controller)
		}

		if err := os.Rename(furl, upload_dir+"img/"+fsname); err != nil {
			beego.Error(err)

			flash.Error(T("resize_bad"))
			flash.Store(&c.Controller)
			return
		}

		d.Img = sql.NullString{upload_dir + "img/" + fsname, true}
	}

	var is_new bool

	if d.Id == 0 {
		is_new = true
		v, err = DM.Insert(&d)
	} else {
		v, err = DM.Update(&d)
	}

	if err != nil {
		beego.Error("Model insert error: ", err)
		flash.Error(T("internal"))
		flash.Store(&c.Controller)
		return
	}

	if v.HasErrors() {
		M.ExpandFormErrors(&v, c.Data)
		return
	}

	if is_new {
		tmpl := c.AddMail(d.Name, d.Id)
		c.SendMails(tmpl)

		c.SendTgNotifs(fmt.Sprintf("New Candidate Added:  %s %s  %s:%d/adm/candidates/%d", d.Name,
			d.LName,
			c.Ctx.Input.Site(),
			c.Ctx.Input.Port(),
			d.Id))
	}

	flash.Notice(T("cand_created", map[string]interface{}{"Name": d.Name, "Lname": d.LName}))
	flash.Store(&c.Controller)

	c.Redirect(beego.UrlFor("Cand.Index"), 302)
}

//add new candidate
func (c *Cand) Add() {

	flash := beego.NewFlash()
	var deps []M.Dep

	cand := M.Cand{}

	err := DM.FindAll(&M.Dep{}, &deps, M.Sf{}, M.Where{})

	if err != nil {
		beego.Error("Department find error: ", err)

		flash.Error(T("internal"))
		flash.Store(&c.Controller)

		return
	}

	c.Data["deps"] = deps
	c.Data["cand"] = cand
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
}

//edit candidate
func (c *Cand) Edit() {

	flash := beego.NewFlash()
	var deps []M.Dep
	var err error

	if err := DM.FindAll(&M.Dep{}, &deps, M.Sf{}, M.Where{}); err != nil {
		beego.Error("Department find error: ", err)

		flash.Error(T("internal"))
		flash.Store(&c.Controller)

		return
	}
	var id int

	if id, err = strconv.Atoi(c.Ctx.Input.Param(":id")); err != nil {
		beego.Error(err)
	}

	var cand M.Cand

	if DM.FindByPk(&cand, id) != nil {
		flash.Notice(T("nocand"))
		flash.Store(&c.Controller)

		c.Redirect(beego.UrlFor("Cand.Index"), 302)
		return
	}

	c.Data["cand"] = cand
	c.Data["deps"] = deps
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
}

//remove candidate
func (c *Cand) Remove() {

	var s string

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		s = T("invalid_param", map[string]interface{}{"Param": "ID"})
	}

	d := M.Cand{}

	err = DM.DeleteByPk(&d, id)

	if err != nil {
		s = T("nocand", 1)
		beego.Error(err)
	} else {
		s = T("cand_deleted")
	}

	c.Data["json"] = RJson{s, err == nil}
}

func (c *Cand) Search() {

	if !c.IsAjax() {
		return
	}

	cands, err := M.Cand{}.Search(c.GetString("q"), 5)

	if err != nil {
		beego.Error(err)
	}

	c.Data["json"] = cands
}

func (c *Cand) AddMail(name string, id int) string {
	c.Layout = "mail-layout.tpl"
	c.TplNames = "mail/candadd.tpl"
	// get link to candidate page
	site := c.Ctx.Input.Site()
	port := strconv.Itoa(c.Ctx.Input.Port())
	// TODO get from url proper candidates view path
	ns := "/adm/candidates"
	c.Data["cand_url"] = site + ":" + port + ns + "/" + strconv.Itoa(id)
	c.Data["cand_id"] = id
	c.Data["cand_name"] = name

	temp, err := c.RenderString()

	if err != nil {
		beego.Error(err)
	}

	return temp
}
