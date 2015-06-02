package models

import (
	"database/sql"
	"testing"
)

func TestAddMail(t *testing.T) {
	var err error

	mail := Mail{}

	mail.ToMail = "ahidoyatov@gmail.com"
	mail.ToName = sql.NullString{"Abdullo", true}
	mail.ToType = "to"
	mail.FromMail = "ahidoyatov@gmail.com"
	mail.FromName = sql.NullString{"hrkb", true}
	mail.Subject = sql.NullString{"notification", true}
	mail.Html = sql.NullString{"<h1>long interestin html</h1>", true}
	mail.Text = sql.NullString{"some text", true}

	_, err = dm.Insert(&mail)
	if err != nil {
		t.Error(err)
	}
}
