package models

type Contact struct {
	Id int `db:"id"`
	Name string `db:"name" form:"name" valid:"Required"`
	Value string `db:"value" form:"value" valid:"Required"`
	Cand int `db:"ref_cand"`
	Active bool `db:"active"`
}

func (m Contact) Table() string {
	return "contacts"
}

func (m *Contact) Reset() {
	*m = Contact{}
}

func (m *Contact) SetActive(b bool) {
	m.Active = b
}

func (m Contact) GetSelf() interface{} {
	return m
}

func (m *Contact) GetId() int {
	return m.Id
}

func (m *Contact) SetId(id int) {
	m.Id = id
}

