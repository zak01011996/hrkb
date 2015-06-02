package models

type Dep struct {
	Id     int    `db:"id" form:"id"`
	Title  string `db:"title" valid:"Required" form:"title"`
	Active bool   `db:"active"`
}

func (u Dep) Table() string {
	return "departments"
}

func (d *Dep) Reset() {
	*d = Dep{}
}

func (d *Dep) SetActive(b bool) {
	d.Active = b
}

func (d Dep) GetSelf() interface{} {
	return d
}

func (d *Dep) GetId() int {
	return d.Id
}

func (d *Dep) SetId(id int) {
	d.Id = id
}

func (d *Dep) Valid(v *ValidMap) {

	andW := W{"Title": d.Title, "Active": true}

	if d.Id > 0 {
		andW["Id!="] = d.Id
	}

	if dm.Find(&Dep{}, Sf{"Id"}, Where{And: andW, Or: W{"Title": d.Title, "Active": false}}) == nil {
		v.SetError("Title", T("valid_depex"))
	}
}
