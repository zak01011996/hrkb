package models

type Role struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Active bool   `db:"active"`
}

func (r Role) Table() string {
	return "roles"
}

func (r *Role) Reset() {
	*r = Role{}
}

func (r *Role) SetActive(b bool) {
	r.Active = b
}

func (r Role) GetSelf() interface{} {
	return r
}

func (r *Role) GetId() int {
	return r.Id
}

func (r *Role) SetId(id int) {
	r.Id = id
}

