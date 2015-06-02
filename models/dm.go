package models

import (
	"github.com/coopernurse/gorp"
)

type DM struct {
	DB *gorp.DbMap
}

func NewDM(db *gorp.DbMap) *DM {
	return &DM{DB: db}
}

func (d *DM) Find(m Model, sf Sf, w Where, p ...*Params) error {
	m.Reset()
	a, o := prepActiveCond(w)
	return find(m, sf, a, o, p...)
}

func (d *DM) FindByPk(m Model, id int) error {
	m.Reset()
	a, o := prepActiveCond(Where{And: W{"Id": id}})
	return find(m, Sf{}, a, o)
}

func (d *DM) FindAll(m Model, dist interface{}, sf Sf, w Where, p ...*Params) error {
	a, o := prepActiveCond(w)
	err := findAll(m, dist, sf, a, o, p...)
	return err
}

func (d *DM) DeleteByPk(m Model, id int) error {
	m.SetId(id)
	m.SetActive(false)

	if _, err := update(m, Sf{"Active"}); err != nil {
		return err
	}

	return nil
}

/**
 Find before delete and fetch result row to model if find no errors
**/
func (d *DM) DeleteByPkWithFetch(m Model, id int) error {
	err := d.FindByPk(m, id)
	if err != nil {
		return err
	}
	return d.DeleteByPk(m, id)
}

func (d *DM) RestoreByPk(m Model, id int) (err error) {

	err = d.Find(m, Sf{}, Where{And: W{"Id": id, "Active": false}})
	if err == nil {
		m.SetActive(true)
		_, err = update(m, Sf{"Active"})
	}

	return err
}

func (d *DM) DeleteByCondition(m Model, w Where) error {
	err := d.Find(m, Sf{}, w)
	if err == nil {
		m.SetActive(false)
		_, err = update(m, Sf{"Active"})
	}

	return err
}

func (d *DM) Insert(m Model, fields ...string) (ValidMap, error) {

	var f Sf
	var err error = nil

	for _, v := range fields {

		if v == "*" {
			break
		}

		f = append(f, v)
	}

	v := Validate(m)
	if !v.HasErrors() {
		err = insert(m, f)
	}
	return v, err
}

func (d *DM) Update(m Model, fields ...string) (ValidMap, error) {

	var f Sf
	var err error = nil

	for _, v := range fields {

		if v == "*" {
			break
		}

		f = append(f, v)
	}

	v := Validate(m)
	if !v.HasErrors() {
		_, err = update(m, f)
	}

	return v, err
}

func (d *DM) Delete(m Model) error {
	_, err := dbmap.Delete(m)
	return err
}

func (d *DM) Count(m Model, w Where) (int64, error) {
	a, o := prepActiveCond(w)
	return count(m, a, o)
}
