package ctrls

import (
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type (
	Trash struct {
		BaseController
	}
	tc struct {
		Id                           int
		Dt, UserName, CandName, Text string
	}
	Mapping struct {
		model  M.Model
		mSlice interface{}
	}
)

var TrashDateFmt string

func getModels() map[string]Mapping {
	var cands []M.Cand
	var deps []M.Dep
	var crits []M.Crit
	var users []M.User
	var comments []M.Comment

	return map[string]Mapping{
		"deps":     Mapping{&M.Dep{}, &deps},
		"cands":    Mapping{&M.Cand{}, &cands},
		"crits":    Mapping{&M.Crit{}, &crits},
		"users":    Mapping{&M.User{}, &users},
		"comments": Mapping{&M.Comment{}, &comments},
	}
}

func (c *Trash) Get() {
	var err error
	models := getModels()
	for k, m := range models {

		err = DM.FindAll(m.model, m.mSlice, M.Sf{}, M.Where{And: M.W{"Active": false}})
		if err != nil {
			beego.Error(err)
			return
		}

		if k == "comments" {
			var t []tc
			var u M.User
			var cand M.Cand
			n := m.mSlice.(*[]M.Comment)
			for _, v := range *n {
				x := tc{Id: v.Id, Text: v.Comment}
				if v.CreatedAt != nil {
					x.Dt = v.CreatedAt.Format(TrashDateFmt)
				}

				if DM.FindByPk(&u, v.User) == nil {
					x.UserName = u.Name
				}

				if DM.FindByPk(&cand, v.Cand) == nil {
					x.CandName = cand.Name + " " + cand.LName
				}

				t = append(t, x)
			}
			c.Data[k] = t
			continue
		}

		c.Data[k] = m.mSlice
	}
}

func (c *Trash) Restore() {
	if !c.IsAjax() {
		return
	}

	jresp := RJson{Msg: T("no_rec")}

	c.Data["json"] = &jresp

	mType := c.Ctx.Input.Param(":type")
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
		return
	}

	models := getModels()
	if _, ok := models[mType]; !ok {
		return
	}

	err = DM.RestoreByPk(models[mType].model, id)

	if err != nil {
		beego.Error(err)
		return
	}
	jresp.Msg = T("restored")
	jresp.Success = true
}
