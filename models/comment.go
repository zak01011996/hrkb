package models

import "time"

type Comment struct {
	Id   int `db:"id"`
	User int `db:"ref_user"`
	Cand int `db:"ref_cand"`

	//can be used as row.CreatedAt.Format("02.01.2006 15:04:05") to evalute to dd.mm.yyyy hh:mm:ss format
	CreatedAt *time.Time `db:"created_dt"`

	Comment string `db:"comment" form:"text" valid:"Required"`
	Active  bool   `db:"active"`
}

func (c Comment) Table() string {
	return "comments"
}

func (c *Comment) Reset() {
	*c = Comment{}
}

func (c *Comment) SetActive(b bool) {
	c.Active = b
}

func (c Comment) GetSelf() interface{} {
	return c
}

func (c *Comment) GetId() int {
	return c.Id
}

func (c *Comment) SetId(id int) {
	c.Id = id
}
