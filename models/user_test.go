package models

import (
	"fmt"
	"testing"
)

func TestPrepActiveCond(t *testing.T) {

	a, o := prepActiveCond(Where{})

	if len(a) == 1 {

		if _, ok := a["Active"]; ok {

			if len(o) > 0 || o != nil {
				t.Error("o is initialized")
			}

		} else {
			t.Error("a not contains Active")
		}

	} else {
		t.Error(fmt.Sprintf("length a=%d", len(a)))
	}

	a, o = prepActiveCond(Where{And: W{"Login": "admin"}})

	if len(a) != 2 || len(o) != 0 {
		t.Error(fmt.Sprintf("a length %d, o length %d", len(a), len(o)))
	}

	a, o = prepActiveCond(Where{And: W{"Login": "admin"}, Or: W{"Active": false}})

	if _, ok := a["Active"]; ok {
		t.Error("a is initialized")
	} else if len(a) != 1 || len(o) != 1 {
		t.Error(fmt.Sprintf("2) a length %d, o length %d", len(a), len(o)))
	}

	a, o = prepActiveCond(Where{And: W{"Login": "admin"}, Or: W{"Login": "manager", "Role": RoleManager}})

	if _, ok := a["Active"]; ok {

		if _, ok := o["Active"]; ok {

			if len(a) != 2 || len(o) != 3 {
				t.Error(fmt.Sprintf("3) a length %d, o length %d", len(a), len(o)))
			}

		} else {
			t.Error("3) o not contains Active")
		}

	} else {
		t.Error("3) a not contains Active")
	}
}

func TestUserCheckPass(t *testing.T) {

	_, err := CheckPass("admin", "1")
	if err != nil {
		t.Error(err, "login=admin & password=1")
	}

	id, err := CheckPass("manager", "wrong_pass")
	if err == nil {
		t.Error(id, "login=manager & password=wrong_pass successfuly authorized!")
	}

}

func TestUserValidation(t *testing.T) {

	u := User{
		Id:       1,
		Login:    "admin",
		Name:     "Double Test",
		Password: "123",
		PassConf: "456",
		Email:    "mirolim777@gmail.com",
		Role:     4,
		Active:   false,
	}

	m := make(map[string]bool, 0)
	m["PassConf"] = false
	m["Role"] = false
	m["Login"] = false

	v := Validate(&u)

	if !v.HasErrors() {
		t.Error("errors expected")
	} else {
		for _, err := range v.Errors {
			if _, ok := m[err.Key]; ok {
				m[err.Key] = true
			}
		}
		for k, v := range m {
			if !v {
				t.Error(k + " skiped!")
			}
		}
	}
	// 701 ID - is test data for admin
	u = User{
		Id:       701,
		Login:    "admin",
		Name:     "Triple Test",
		Password: "123",
		PassConf: "123",
		Role:     RoleManager,
		Email:    "mirolim777@gmail.com",
		Active:   false,
	}

	v = Validate(&u)

	if v.HasErrors() {
		t.Error("no errors expected")
	}

}
