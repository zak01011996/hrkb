package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
	"hrkb/utils"
)

type User struct {
	Id               int            `db:"id" form:"id"`
	Login            string         `db:"login" valid:"Required;AlphaNumeric" form:"login"`
	Name             string         `db:"name" valid:"Required" form:"name"`
	Password         string         `db:"password" valid:"Required" form:"pass"`
	PassConf         string         `db:"-" valid:"Required" form:"passc"`
	Email            string         `db:"mail" form:"mail"`
	GToken           sql.NullString `db:"gitlab_token" form:"gtoken"`
	Role             int            `db:"role" valid:"Required" form:"role"`
	NotifyByMail     bool           `db:"notify_by_mail" form:"notifyMail"`
	NotifyByTelegram bool           `db:"notify_by_telegram" form:"notifyTelegram"`
	Active           bool           `db:"active" form:"active"`
}

const (
	RoleAdmin   = 1
	RoleManager = 2
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u User) Table() string {
	return "users"
}

func (u *User) Reset() {
	*u = User{}
}

func (u *User) SetActive(b bool) {
	u.Active = b
}

func (u User) GetSelf() interface{} {
	return u
}

func (u *User) GetId() int {
	return u.Id
}

func (u *User) SetId(id int) {
	u.Id = id
}

func CheckPass(login, pass string) (int, error) {

	var user User

	err := dm.Find(&user, Sf{"Id", "Role", "Password"}, Where{And: W{"Login": login}})
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	}

	if err == nil {
		return user.Id, nil
	}

	return 0, err
}

/*
	additional validation rules
	automaticaly called if base validation rules passed
*/

func (u *User) Valid(v *ValidMap) {

	if u.Password != u.PassConf {
		v.SetError("PassConf", T("valid_passmatch"))
	}

	if dm.Find(&Role{}, Sf{}, Where{And: W{"Id": u.Role}}) != nil {
		v.SetError("Role", T("valid_role"))
	}

	andW := W{"Login": u.Login, "Active": true}

	if u.Id > 0 {
		andW["Id!="] = u.Id
	}

	if dm.Find(&User{}, Sf{"Id"}, Where{And: andW, Or: W{"Login": u.Login, "Active": false}}) == nil {
		v.SetError("Login", T("valid_userex"))
	}

	if !v.HasErrors() {
		if pass, err := utils.HashPass(u.Password); err != nil {
			v.SetError("Password", T("valid_hash"))
		} else {
			u.Password = pass
		}
	}
}
