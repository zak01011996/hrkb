package ctrls

import (
	"database/sql"
	"strconv"

	"github.com/astaxie/beego"
	M "hrkb/models"
)

type Rat struct {
	BaseController
}

// get concrete candidate grouped ratings via ajax
// get param :id is candidate id
func (c *Rat) Get() {
	if !c.IsAjax() {
		return
	}

	var err error
	var id int
	var ratings []M.GroupedRating

	f := &M.Rat{}
	id, err = strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
	}
	ratings, err = f.GetGroupedRatings(id)

	if err != nil {
		beego.Error(err)
	}

	c.Data["json"] = struct {
		Ratings []M.GroupedRating
	}{ratings}
}

// get concrete candidate detailed ratings via ajax
// get param :id is candidate id
func (c *Rat) GetDetailed() {
	if !c.IsAjax() {
		return
	}

	var err error
	var id int
	var ratings []M.DetailedRating

	f := &M.Rat{}
	id, err = strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
	}
	ratings, err = f.GetDetailedRatings(id)

	if err != nil {
		beego.Error(err)
	}

	c.Data["json"] = struct {
		Ratings []M.DetailedRating
	}{ratings}
}

// get concrete concrete user ratings via ajax
// get param :id is candidate id
func (c *Rat) GetUserRatings() {
	if !c.IsAjax() {
		return
	}

	var err error
	var id int
	var ratings []M.UserRating

	f := &M.Rat{}
	id, err = strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		beego.Error(err)
	}

	i := c.GetSession("uid")
	if i == nil {
		beego.Error(err)
		return
	}

	ratings, err = f.GetUserRatings(id, i.(int))

	if err != nil {
		beego.Error(err)
	}

	c.Data["json"] = struct {
		Ratings []M.UserRating
	}{ratings}
}

// save concrete candidate rating
func (c *Rat) Post() {

	if !c.IsAjax() {
		return
	}

	var err error
	var s string
	var v M.ValidMap
	var ok bool

	r := M.Rat{Active: true}

	err = c.ParseForm(&r)

	uid := c.GetSession("uid")
	if uid == nil {
		c.Data["json"] = RJson{T("auth_is_off"), ok}
		return
	}

	fr := M.Rat{}
	err = DM.Find(&fr, M.Sf{}, M.Where{And: M.W{"Cand": r.Cand, "User": uid.(int), "Crit": r.Crit}})
	if err != nil && err != sql.ErrNoRows {
		c.Data["json"] = RJson{T("internal"), ok}
		beego.Warn(err)
		return
	}

	if fr.Id > 0 {
		c.Data["json"] = RJson{T("rat_already"), ok}
		return
	}

	r.User = uid.(int)

	if err == sql.ErrNoRows {
		v, err = DM.Insert(&r)
	}

	if err != nil {
		s = T("save_error")
		beego.Error(err)
	}

	if !v.HasErrors() {
		s = T("rat_added")
		ok = true
	} else {
		s = T("validation_error")
	}

	c.Data["json"] = RJson{s, ok}
}

func (c *Rat) Remove() {

	s := T("rat_deleted")

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		s = T("invalid_param", map[string]interface{}{"Param": "ID"})
	}

	uid := c.GetSession("uid")
	if uid == nil {
		c.Data["json"] = RJson{T("auth_is_off"), err == nil}
		return
	}

	r := M.Rat{}

	err = DM.DeleteByCondition(&r, M.Where{And: M.W{"User": uid.(int), "Id": id}})

	if err != nil {
		s = T("rat_not_found")
	}

	c.Data["json"] = RJson{s, err == nil}
}
