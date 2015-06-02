package models

import "hrkb/utils"

type Profile struct {
	Id               int    `db:"id"`
	Name             string `valid:"Required" form:"name"`
	Email            string `valid:"Email" db:"mail" form:"email"`
	Password         string `form:"pass"`
	PassConf         string `form:"passc"`
	NotifyByMail     bool   `form:"notifyMail"`
	NotifyByTelegram bool   `form:"notifyTelegram"`
}

func (m *Profile) Valid(v *ValidMap) {

	if m.Password == "" {
		return
	}

	if m.Password != m.PassConf {
		v.SetError("PassConf", T("valid_passmatch"))
		return
	}

	pass, err := utils.HashPass(m.Password)

	if err != nil {
		v.SetError("Password", T("valid_hash"))
		return
	}

	m.Password = pass
}

func (m Profile) Table() string {
	return "users"
}

func (m *Profile) Reset() {
	*m = Profile{}
}

func (m *Profile) SetActive(b bool) {

}

func (m Profile) GetSelf() interface{} {
	return m
}

func (m *Profile) GetId() int {
	return m.Id
}

func (m *Profile) SetId(id int) {
	m.Id = id
}
