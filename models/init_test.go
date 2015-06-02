package models

import (
	"fmt"
	"testing"

	"hrkb/conf"
	"hrkb/utils"
)

func TestDbOpen(t *testing.T) {
	if err := DbOpen(conf.Database{Driver: "???"}); err == nil {
		t.Error("DbOpen opened with undefined driver")
	}
}

func TestSplitOp(t *testing.T) {
	a := [][3]string{
		{"Login>", "Login", ">"},
		{"Age!=", "Age", "!="},
		{"Name", "Name", "="},
		{"Name<", "Name", "<"},
	}
	for i, j := 0, len(a); i < j; i += 1 {
		s, op := splitOp(a[i][0])
		if s != a[i][1] || op != a[i][2] {
			t.Error(fmt.Sprintf("%s failed. Actual results is %s %s. Expected: %s %s", a[i][0], s, op, a[i][1], a[i][2]))
		}
	}
}

func TestConcatConds(t *testing.T) {
	mock := []struct {
		op   string
		args []string
		ret  string
	}{
		{"AND", []string{"1", "2", "3"}, "1 AND 2 AND 3"},
		{"OR", []string{"1", "2"}, "1 OR 2"},
		{"OR", []string{"", "", ""}, ""},
	}

	for k, v := range mock {
		if s := concatConds(v.op, v.args...); s != v.ret {
			t.Error(fmt.Sprintf("%d) %s expected actual %s", k, v.ret, s))
		}
	}
}

func TestPrepSelect(t *testing.T) {

	mock := []struct {
		m         Model
		sf        Sf
		andW, orW W
		p         *Params
		result    int
	}{
		{&User{}, Sf{"Id", "Login", "Password"}, W{"Id>": 10, "Login!=": "admin"}, W{}, NewParams(Params{Limit: 10}), 3},

		{&User{}, Sf{}, W{}, W{}, nil, 0},

		{&User{}, Sf{}, W{"Id>": 10, "Login!=": "admin"}, W{"Id": 2, "Login": "manager"}, nil, 4},
	}

	var d sdictype

	for i, j := 0, len(mock); i < j; i += 1 {

		table := mock[i].m.Table()
		if _, ok := dic[table]; ok {
			d = dic[table]
		}

		sql, pv := prepSelect(table, d, mock[i].sf, mock[i].andW, mock[i].orW, mock[i].p)
		if len(pv) != mock[i].result {
			t.Error(fmt.Sprintf("#%d %s %v", i, sql, pv))
		}
	}

}

func TestFind(t *testing.T) {

	u := User{}

	if err := find(&u, Sf{"Id", "Login"}, W{"Id>": 0}, W{"Active": true}, NewParams(Params{Sort: "Id DESC"})); err != nil {
		t.Error(err)
	}

	if err := find(&u, Sf{}, W{"Id!=": -1}, W{}); err != nil {
		t.Error(err)
	}
}

func TestFindAll(t *testing.T) {

	var a, b []User

	if err := findAll(&User{}, &a, Sf{}, W{}, W{}, NewParams(Params{Sort: "Login ASC"})); err != nil {
		t.Error(err)
	}

	if err := findAll(&User{}, &b, Sf{"Login"}, W{}, W{}, NewParams(Params{Limit: 2})); err != nil || len(b) != 2 {
		t.Error(err, fmt.Sprintf("length b=%d", len(b)))
	}

}

func TestParams_BuildSql(t *testing.T) {

	p := NewParams(Params{Sort: "Id ASC, Login DESC", Offset: 5})

	if s := p.buildSql(dic["users"], &pvalues{}); s != " ORDER BY id ASC, login DESC OFFSET $1" {
		t.Error(s)
	}
}

func TestInsert(t *testing.T) {
	u := User{Login: "Bruce Willis", Name: "TTTT", Password: "123", Role: 2, Email: "ahidoyatov@gmail.com"}
	err := insert(&u, Sf{"Login", "Name", "Password", "Role", "Email"})
	if err != nil {
		t.Error(err)
	} else if u.Login != "Bruce Willis" || u.Role != 2 || u.Active != true {
		t.Error("ivalid return values", u)
	}

	err = insert(&u, Sf{"L0gin"})
	if err == nil {
		t.Error("no insert excepted")
	}
}

func TestUpdate(t *testing.T) {
	mock := []struct {
		m  Model
		sf Sf
		er bool
	}{
		{&User{Id: 702, Login: "gomer simpson"}, Sf{"Longin"}, true},
		{&User{Id: 702, Login: "gomer simpson"}, Sf{"Id"}, true},
		{&User{Id: 702, Login: "gomer simpson", Active: true}, Sf{"Login"}, false},
		{&User{Id: 702, Login: "gomer simpson", Active: true}, Sf{"Login", "Active"}, false},
	}
	for k, v := range mock {
		n, err := update(v.m, v.sf)
		if !v.er && err != nil {
			t.Error(err)
		}
		if err == nil && n != 1 {
			t.Error(k, fmt.Sprintf("expected 1 row update. actual %d", n))
		}
		if err == nil && n == 1 {
			u := User{}
			err := find(&u, Sf{"Login"}, W{"Id": v.m.(*User).Id}, W{})
			if err != nil {
				t.Error(err)
			} else if utils.IndexOfStr(v.sf, "Login") != -1 && u.Login != v.m.(*User).Login {
				t.Error(fmt.Sprintf("%d updated to %s", k, u.Login))
			}
		}
	}
}

func TestTransactions(t *testing.T) {

	m := []struct {
		k string
		u *User
	}{
		{"Begin", nil},
		{"Insert", &User{Login: "Cameron Diaz", Password: "Cameron Diaz", Role: 2}},
		{"Update", &User{Login: "Cameron Diaz 1", Password: "Some", Role: 2}},
		{"Commit", nil},
		{"Begin", nil},
		{"Delete", &User{}},
		{"Rollback", nil},
		{"Begin", nil},
		{"Delete", &User{}},
		{"Commit", nil},
	}

	err := make(chan error)
	var tr *Transaction
	var id int

	for _, v := range m {
		go func(k string, u *User) {

			var s error

			switch k {
			case "Begin":
				tr, s = Begin()
				err <- s
			case "Insert":
				s = tr.Insert(u)
				// fmt.Printf("%+v \n", u)
				if s == nil {
					id = u.Id
				}
				err <- s
			case "Update":
				u.Id = id
				_, s = tr.Update(u)
				err <- s
			case "Delete":
				u.Id = id
				_, s = tr.Delete(u)
				err <- s
			case "Commit":
				err <- tr.Commit()
			case "Rollback":
				err <- tr.Rollback()
			default:
				err <- nil
			}

		}(v.k, v.u)

		if s := <-err; s != nil {
			t.Error(s)
			break
		}
	}

}
