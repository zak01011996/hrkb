package models

import "time"

type File struct {
	Id int `db:"id"`
	Cand int `db:"ref_cand"`
	User int `db:"ref_user"`
	Url string `db:"url"`
	Created time.Time `db:"created"`
	Name string `db:"name"`
	Mime string `db:"mime"`
	Size int64 `db:"fsize"`
	Active bool `db:"active"`
}

func (m File) Table() string {
	return "files"
}

func (m *File) Reset() {
	*m = File{}
}

func (m *File) SetActive(b bool) {
	m.Active = b
}

func (m File) GetSelf() interface{} {
	return m
}

func (m *File) GetId() int {
	return m.Id
}

func (m *File) SetId(id int) {
	m.Id = id
}

