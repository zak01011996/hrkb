package ctrls

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/nicksnyder/go-i18n/i18n"
	M "hrkb/models"
)

type Lang struct {
	BaseController
}

func (c *Lang) Index() {
	var langs []M.Lang
	err := DM.FindAll(&M.Lang{}, &langs, M.Sf{"Id", "Code", "IsDefault"}, M.Where{})
	if err != nil {
		beego.Error(err)
	}
	c.Data["langs"] = langs
	c.Data["langsSup"] = M.LangsSupports
	c.Data["title"] = T("lang", 2)
	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["code"] = ""
}

func (c *Lang) Remove() {

	s := T("lang_not_found")
	d := M.Lang{}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.DeleteByPkWithFetch(&d, id)
	}

	if err == nil {
		s = T("lang_removed")
		err = c.cleanLang(d.Code)
	}

	if err != nil {
		beego.Error(err)
	}

	if c.IsAjax() {
		c.Data["json"] = RJson{s, err == nil}
		return
	}

	flash := beego.NewFlash()
	flash.Notice(s)
	flash.Store(&c.Controller)

	c.Redirect(beego.UrlFor("Lang.Index"), 302)

}

func (c *Lang) Upload() {

	r := RJson{T("lang_not_found"), false}
	l := &M.Lang{}
	c.Data["json"] = &r
	c.EnableRender = false

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = DM.FindByPk(l, id)
	}

	if err == nil {
		err, r.Msg = c.applyFile(l.Code)
	}

	if err != nil {
		beego.Error(err)
		c.ServeJson()
		return
	}

	r.Msg = T("lang_fileok")
	r.Success = true
	c.ServeJson()
}

func (c *Lang) Download() {
	var l M.Lang
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err == nil {
		err = DM.FindByPk(&l, id)
	}

	if err != nil {
		beego.Error(err)
		c.Redirect(beego.UrlFor("Lang.Index"), 302)
	}

	c.Ctx.Output.Download(c.langFileName("lang::folder", l.Code))
}

func (c *Lang) Default() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err == nil {
		err = M.Lang{}.SetDefault(id)
	}

	if err != nil {
		beego.Error(err)
	}

	c.Redirect(beego.UrlFor("Lang.Index"), 302)
}

func (c *Lang) Add() {

	l := &M.Lang{}
	err := c.ParseForm(l)
	if err != nil {
		beego.Error(err)
		c.Redirect(beego.UrlFor("Lang.Index"), 302)
		return
	}

	v := M.Validate(l)

	if v.HasErrors() {
		c.Index()
		c.Data["code"] = l.Code
		c.TplNames = "lang/index.tpl"
		M.ExpandFormErrors(&v, c.Data)
		return
	}

	err, msg := c.applyFile(l.Code)
	flash := beego.NewFlash()
	if err != nil {

		if DM.Find(l, M.Sf{"Code"}, M.Where{And: M.W{"Code": l.Code}}) != nil {

			err := i18n.LoadTranslationFile(c.langFileName("lang::folder", l.Code))
			if err != nil {
				beego.Error(err)
			}
		}

		beego.Error(err)
		flash.Notice(msg)
		flash.Store(&c.Controller)
		c.Redirect(beego.UrlFor("Lang.Index"), 302)
		return
	}

	code := l.Code

	if DM.Find(l, M.Sf{"Id", "Code"}, M.Where{And: M.W{"Code": l.Code, "Active": true}, Or: M.W{"Code": l.Code, "Active": false}}) == nil {
		if !l.Active {
			l.Active = true
			_, err = DM.Update(l, "Active")
		}
	} else {
		l.Code = code
		v, err = DM.Insert(l, "Code")
	}

	if err != nil {
		beego.Error(err)
		if err := c.cleanLang(l.Code); err != nil {
			beego.Error(err)
		}
		flash.Notice(T("internal"))
		flash.Store(&c.Controller)
	}

	c.Redirect(beego.UrlFor("Lang.Index"), 302)

}

/**
Clean language buffer with sending empty file
For example it can be needed on deleting language or more
**/
func (c *Lang) cleanLang(code string) error {
	s := c.langFileName("tmp_dir", code)
	err := ioutil.WriteFile(s, []byte("[]"), 0600)
	if err != nil {
		return err
	}
	defer os.Remove(s)
	return i18n.LoadTranslationFile(s)
}

/**
Try to receive sended via post file and apply it to language by param code
returns error if occured by validation or adding to translations
**/
func (c *Lang) applyFile(code string) (error, string) {

	_, header, err := c.GetFile("file")
	if err != nil {
		return err, T("lang_nofile")
	}

	a := strings.Split(header.Filename, ".")
	if l := len(a); l < 2 || strings.ToLower(a[l-1]) != "json" {
		return errors.New("client validation by file type hacked"), T("lang_badfile")
	}

	s := c.langFileName("tmp_dir", code)

	err = c.SaveToFile("file", s)
	if err != nil {
		return err, T("internal")
	}

	err = i18n.LoadTranslationFile(s)
	defer os.Remove(s)
	if err != nil {
		return err, T("internal")
	}

	s2 := c.langFileName("lang::folder", code)
	err = c.SaveToFile("file", s2)
	if err == nil {
		err = i18n.LoadTranslationFile(s2)
	}

	if err != nil {
		return err, T("internal")
	}

	return nil, ""
}

func (c *Lang) langFileName(cfgFolderKey, code string) string {
	return beego.AppConfig.String(cfgFolderKey) + code + ".all.json"
}
