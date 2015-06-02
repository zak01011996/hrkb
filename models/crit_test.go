package models

import "testing"

func TestGetCritGroup(t *testing.T) {
	c := &Crit{}
	groups, err := c.GetGroupedCrits()
	if err != nil {
		t.Error(err)
	}

	if len(groups) <= 0 {
		t.Error("No data for test")
	}
	for _, v := range groups {

		switch v.DepId {
		case 700:
			if v.Department != "Test department" {
				t.Error("Department not equals", v.DepId)
			}
			if len(v.Criterias) != 0 {
				t.Error("Structure error on department ID ", v.DepId)
			}
		case 701:
			if v.Department != "Test department for criteria 1" {
				t.Error("Department not equals 701")
			}

			if len(v.Criterias) <= 0 {
				t.Error("Structure error on department ID ", v.DepId)
			}

			for _, crit := range v.Criterias {
				if !intInSlice(crit.Id, []int{700, 701, 702}) {
					t.Error("Wrong criterias in department ")
				}
			}
		case 702:
			if v.Department != "Test department for criteria 2" {
				t.Error("Department not equals ", v.DepId)
			}

			if len(v.Criterias) <= 0 {
				t.Error("Structure error on department ID ", v.DepId)
			}

			for _, crit := range v.Criterias {
				if !intInSlice(crit.Id, []int{703, 704, 705}) {
					t.Error("Wrong criterias in department ")
				}
			}

		}
	}
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
