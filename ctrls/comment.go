package ctrls

import (
	"html"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Comments struct {
	BaseController
}

type tplComment struct {
	Id                         int
	User                       int `json:"-"`
	Msg, Comment, Author, Date string
}

var CandDateFmt string

func (tc *tplComment) SetAuthor(user_id int) {
	u := M.User{}
	if DM.FindByPk(&u, user_id) == nil {
		tc.Author = u.Name
	}
}

func (tc *tplComment) SetDate(t *time.Time) {
	tc.Date = t.Format(CandDateFmt)
}

func (c *Comments) Post() {

	if !c.IsAjax() {
		return
	}

	var v M.ValidMap
	var cm M.Comment
	var err error

	cr := tplComment{Msg: T("internal")}

	c.Data["json"] = &cr

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err == nil {
		err = c.ParseForm(&cm)
	}

	if err != nil {
		beego.Error(err)
		return
	}

	i := c.GetSession("uid")
	if i == nil {
		cr.Msg = T("auth_is_off")
		return
	}

	cm.User = i.(int) //comment author
	cm.Cand = id      //candidate Id

	cm.Comment = html.EscapeString(cm.Comment)

	v, err = DM.Insert(&cm, "User", "Cand", "Comment")

	if err != nil {
		beego.Error(err)
		return
	}

	if !v.HasErrors() { //if no has errors returns inserted data to client
		cr.Id = cm.Id
		cr.Msg = "Ok"
		cr.Comment = cm.Comment
		cr.SetAuthor(cm.User)
		cr.SetDate(cm.CreatedAt)
		return
	}

	M.ExpandFormErrors(&v, c.Data)

	if val, ok := c.Data["errComment"]; ok {
		cr.Msg = val.(string) //return error by field comment
	}
}

func (c *Comments) Edit() {

	cr := struct {
		Text string `form:"text" json:"text"`
		Ok   bool   `json:"success"`
		Msg  string `json:"error"`
		Dt   string `json:"dt"`
	}{Msg: T("internal")}

	c.Data["json"] = &cr

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err == nil {
		err = c.ParseForm(&cr)
	}
	if err != nil {
		beego.Error(err)
		cr.Text = ""
		return
	}

	i := c.GetSession("uid")
	if i == nil {
		cr.Msg = T("auth_is_off")
		cr.Text = ""
		return
	}

	var cm M.Comment

	err = DM.Find(&cm, M.Sf{"Id", "CreatedAt"}, M.Where{And: M.W{"Id": id, "User": i.(int)}})
	if err != nil {
		cr.Msg = T("comment_not_found")
		cr.Text = ""
		return
	}

	cm.Comment = html.EscapeString(cr.Text)
	now := time.Now()
	*cm.CreatedAt = now
	v, err := DM.Update(&cm, "Comment", "CreatedAt")

	if v.HasErrors() {
		for _, err := range v.Errors {
			cr.Msg = err.Message
			break
		}
		cr.Text = ""
		return
	}

	if err != nil {
		beego.Error(err)
		cr.Text = ""
		return
	}

	cr.Text = cm.Comment
	cr.Ok = true
	cr.Dt = now.Format(CandDateFmt)
	cr.Msg = ""
}

func (c *Comments) Remove() {

	r := RJson{Msg: T("internal")}
	c.Data["json"] = &r

	i := c.GetSession("role")
	if i == nil || i.(int) != M.RoleAdmin {
		r.Msg = T("restrict_access")
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		return
	}

	err = DM.DeleteByPk(&M.Comment{}, id)
	if err != nil {
		beego.Error(err)
		return
	}

	r.Msg = ""
	r.Success = true
}
