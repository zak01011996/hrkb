package models

import (
	"database/sql"
	"time"
)

// Mailidate main structure
type Mail struct {
	Id       int            `db:"id"`
	ToMail   string         `db:"to_mail" valid:"Required"`
	ToName   sql.NullString `db:"to_name"`
	ToType   string         `db:"to_type" valid:"Required"`
	FromMail string         `db:"from_mail" valid:"Required"`
	FromName sql.NullString `db:"from_name"`
	Subject  sql.NullString `db:"subject"`
	Html     sql.NullString `db:"html"`
	Text     sql.NullString `db:"text"`
	Try      int            `db:"try"`
	Created  time.Time      `db:"created"`
	Updated  time.Time      `db:"updated"`
	Status   bool           `db:"status"`
	Active   bool           `db:"active"`
}

func (m Mail) Table() string {
	return "mails"
}

func (m *Mail) Reset() {
	*m = Mail{}
}

func (m *Mail) SetActive(b bool) {
	m.Active = b
}

func (m Mail) GetSelf() interface{} {
	return m
}

func (m *Mail) GetId() int {
	return m.Id
}

func (m *Mail) SetId(id int) {
	m.Id = id
}
