package ctrls

import (
	"html/template"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Download struct {
	BaseController
	dest string
}

func (c *Download) Index() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
	}

	var files []M.File

	if err := DM.FindAll(&M.File{}, &files, M.Sf{}, M.Where{And: M.W{"Cand": id}}); err != nil {
		beego.Error(err)
	}

	c.Data["xsrfdata"] = template.HTML(c.XsrfFormHtml())
	c.Data["files"] = files
	c.Data["title"] = "Files"
}

//Download file with any type
func (c *Download) Get() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error()
	}

	c.dest = beego.AppConfig.String("upload_dir") + strconv.Itoa(id) + "/"

	file := M.File{}

	if err := DM.Find(&file, M.Sf{}, M.Where{And: M.W{"Id": id}}); err != nil {
		beego.Error(err)
	}

	c.Ctx.Output.Download(file.Url)
}

func (c *Download) Remove() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error()
	}

	c.dest = beego.AppConfig.String("upload_dir") + strconv.Itoa(id) + "/"

	file := M.File{}

	if err := DM.Find(&file, M.Sf{}, M.Where{And: M.W{"Id": id}}); err != nil {
		beego.Error(err)
		c.Data["json"] = RJson{T("file_not_found"), false}
		return
	}

	uid := c.GetSession("uid").(int)
	role := c.GetSession("role").(int)

	if !(uid == file.User || role == M.RoleAdmin) {
		c.Ctx.Output.SetStatus(403)

		c.Data["json"] = RJson{T("file_adenied"), false}

		beego.Error("User is not a file owner")
		return
	}

	if err := os.Remove(file.Url); err != nil {
		beego.Error(file.Url, err)

		c.Data["json"] = RJson{T("file_del_error"), false}
		return
	}

	if err := DM.DeleteByPk(&file, id); err != nil {
		beego.Error(err)

		c.Data["json"] = RJson{T("internal"), false}
		return
	}

	c.Data["json"] = RJson{T("file_del_ok"), true}
}
