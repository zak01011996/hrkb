package models

import "testing"

func TestInsertContact(t *testing.T) {
	var err error
	var v ValidMap
	errs := make(map[interface{}]interface{}, 10)

	cont := Contact{Cand: 700, Active: true, Name: "mail", Value: "some@gmail.com" }

	v, err = dm.Insert(&cont)

	if err != nil {
		t.Error(err)
	}

	if v.HasErrors(){
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors unexpected %v", errs)
	}

	cont = Contact{Cand: 700, Active: true}

	v, err = dm.Insert(&cont)

	if err != nil {
		t.Error(err)
	}

	if !v.HasErrors(){
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors unexpected %v", errs)
	}
}

func TestUpdateContact(t *testing.T) {
	var err error
	var v ValidMap
	errs := make(map[interface{}]interface{}, 10)

	cont := Contact{}

	if err := dm.Find(&cont, Sf{}, Where{And: W{"Id": 700}}); err != nil {
		t.Error(err)
	}

	cont.Name = "another"
	cont.Value = "contact"

	v, err = dm.Update(&cont)

	if err != nil {
		t.Error(err)
	}

	if v.HasErrors(){
		ExpandFormErrors(&v, errs)
		t.Errorf("Validation Errors unexpected %v", errs)
	}

	cont.Name = ""
	cont.Value = ""

	v, err = dm.Update(&cont)
	if err != nil {
		t.Error(err)
	}

	if !v.HasErrors(){
		t.Errorf("Validation Errors expected %v", errs)
	}
}

