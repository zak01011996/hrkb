package models

import (
	"fmt"
	"testing"
)

func TestDMFind(t *testing.T) {

	d := &Dep{}
	if err := dm.Find(d, Sf{"Id", "Title", "Active"}, Where{And: W{"Title": "Test department"}}); err != nil {
		t.Error(err)
	}
}

func TestDMFindByPk(t *testing.T) {
	d := &Dep{}
	// 700 ID - Test depertment
	if err := dm.FindByPk(d, 700); err != nil {
		t.Error(err)
	}
}

func TestDMFindAll(t *testing.T) {

	d := &Dep{}
	deps := []Dep{}
	err := dm.FindAll(d, &deps, Sf{}, Where{Or: W{"Active": true}}, NewParams(Params{Limit: 10}))
	if err != nil {
		t.Error(err)
	} else if len(deps) < 1 {
		t.Error("length=0")
	}

}

func TestDMDeleteByPk(t *testing.T) {
	d := &Dep{}
	if err := dm.FindByPk(d, 700); err != nil {
		t.Error(err)
	} else {
		if d.Active {
			if err := dm.DeleteByPk(d, d.Id); err != nil {
				t.Error(err)
			} else {
				d.Active = true
				_, err := dm.Update(d, "Active")
				if err != nil {
					t.Error(err)
					t.Error("department by id=1 need set active on true")
				}
			}
		}
	}
}

var insertId int

func TestDMInsert(t *testing.T) {

	u := &User{Login: "Jesica", Name: "TEST SHMEST", Password: "123", PassConf: "123", Role: RoleAdmin, Email: "ahidoyatov@gmail.com"}

	v, err := dm.Insert(u, "Login", "Name", "Password", "Role", "Email")
	if err != nil || v.HasErrors() {
		t.Error(err, v.ErrorsMap)
		return
	}

	if err := dm.FindByPk(u, u.Id); err == nil && u.Login == "Jesica" {
		insertId = u.Id
	} else {
		t.Error("Inserted user by login Jesica not found", err, u.Login)
	}

}

func TestDMUpdate(t *testing.T) {
	if insertId > 0 {

		u := &User{Id: insertId, Login: "Jesica2222", Name: "YO!", Password: "Alba", PassConf: "Alba", Role: RoleManager}

		v, err := dm.Update(u, "Role")

		if err != nil || v.HasErrors() {
			t.Error(err, v.ErrorsMap)
			return
		}

		if err := dm.FindByPk(u, u.Id); err != nil || u.Role != RoleManager || u.Login != "Jesica" {
			t.Error("Updated user by login Jesica not updated role to manager", err, u)
		}

	}
}

func TestDMDelete(t *testing.T) {
	if insertId > 0 {

		err := dm.Delete(&User{Id: insertId})

		if err != nil {
			t.Error(err)
			return
		}

		if err := dm.FindByPk(&User{}, insertId); err == nil {
			t.Error("Jesica exists in db")
		}

	}
}

func TestDMCount(t *testing.T) {
	c, err := dm.Count(&User{}, Where{})
	if err != nil || c == 0 {
		t.Error(err, c)
	}

	c, err = dm.Count(&User{Login: "admin"}, Where{And: W{"Id": -1}})
	if err != nil || c > 0 {
		t.Error(err, "count is > 0")
	}

	c, err = dm.Count(&User{}, Where{})
	if err != nil || c != 3 {
		t.Error(fmt.Sprintf("3 excepted actual %d", c))
	}
}
