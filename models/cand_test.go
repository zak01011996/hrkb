package models

import "testing"

func TestAddCand(t *testing.T) {
	var err error
	var v ValidMap
	errs := make(map[interface{}]interface{}, 10)

	cand := Cand{Name: "Abdullo", LName: "Xidoyatov", FName: "Some name", Phone: "1233454", Email: "examp@mail.com", Address: "some boring address", Married: false, Dep: 700, Salary: 1000, Currency: "$", Active: true}

	v, err = dm.Insert(&cand)

	if err != nil {
		t.Error(err)
	}

	if v.HasErrors() {
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors unexpected %v", errs)
	}

	cand = Cand{Name: "", LName: "", FName: ""}

	v, err = dm.Insert(&cand)

	if err != nil {
		t.Error(err)
	}

	if !v.HasErrors() {
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors expected %v", v.Errors)
	}
}

func TestUpdateCand(t *testing.T) {
	var err error
	var v ValidMap
	errs := make(map[interface{}]interface{}, 10)

	cand := Cand{}

	dm.Find(&cand, Sf{}, Where{And: W{"Id": 700}})

	cand.Name = "Another Name"
	cand.LName = "Another Lname"
	cand.FName = "Another Fname"
	cand.Salary = 1
	cand.Currency = "sum"

	v, err = dm.Update(&cand)

	if err != nil {
		t.Error(err)
	}

	if v.HasErrors() {
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors unexpected %v", errs)
	}

	cand.Name = ""
	cand.LName = ""
	cand.FName = ""
	cand.Salary = 0
	cand.Currency = ""

	v, err = dm.Update(&cand)

	if err != nil {
		t.Error(err)
	}

	if !v.HasErrors() {
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors expected %v", errs)
	}
}
