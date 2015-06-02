package models

type Crit struct {
	Id     int    `db:"id" form:"id"`
	Title  string `db:"title" valid:"Required" form:"title"`
	Dep    int    `db:"ref_dep" valid:"Required" form:"dep"`
	Active bool   `db:"active"`
}

type CritGroup struct {
	DepId      int
	Department string
	Criterias  []Crit
}

func (u Crit) Table() string {
	return "criteria"
}

func (c *Crit) Reset() {
	*c = Crit{}
}

func (c *Crit) SetActive(b bool) {
	c.Active = b
}

func (c Crit) GetSelf() interface{} {
	return c
}

func (c *Crit) GetId() int {
	return c.Id
}

func (c *Crit) SetId(id int) {
	c.Id = id
}

func (c *Crit) GetGroupedCrits() (ret []CritGroup, err error) {
	d := &Dep{}
	deps := []Dep{}
	err = dm.FindAll(d, &deps, Sf{}, Where{}, NewParams(Params{Sort: "Title ASC"}))

	if err != nil {
		return
	}

	if len(deps) > 0 {
		crits := []Crit{}
		err = dm.FindAll(c, &crits, Sf{}, Where{}, NewParams(Params{Sort: "Title ASC"}))
		if err != nil {
			return
		}
		for _, dp := range deps {
			cts := []Crit{}
			for _, ct := range crits {
				if dp.Id == ct.Dep {
					cts = append(cts, ct)
				}
			}
			if len(cts) > 0 {
				ret = append(ret, CritGroup{dp.Id, dp.Title, cts})
			}
		}
	}

	return
}
