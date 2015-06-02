package models

import (
	"database/sql"
	"strings"
)

// Candidate main structure
type Cand struct {
	Id       int            `db:"id" form:"id"`
	Img      sql.NullString `db:"img" form:"img"`
	Name     string         `db:"name" form:"name" valid:"Required"`
	LName    string         `db:"lname" form:"lname" valid:"Required"`
	FName    string         `db:"fname" form:"fname" valid:"Required"`
	Phone    string         `db:"phone" form:"phone" valid:"Numeric"`
	Email    string         `db:"email" form:"email" valid:"Email"`
	Note     sql.NullString `db:"note" form:"note"`
	Address  string         `db:"addr" form:"addr" valid:Required`
	Married  bool           `db:"married" form:"married"`
	Dep      int            `db:"ref_dep" form:"depId"`
	Salary   float64        `db:"salary" form:"salary" valid:"Required"`
	Currency string         `db:"currency" form:"currency" valid:"Required"`
	Active   bool           `db:"active"`
}

//RCands separated for parsing custom sql search form db and store it in this struct and handle
type RCands struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	LName    string `db:"lname"`
	Dep      string `db:"title"`
	Tel      string `db:"phone"`
	Mail     string `db:"email"`
	Addr     string `db:"addr"`
	Salary   string `db:"salary"`
	Currency string `db:"currency"`
}

func (u Cand) Table() string {
	return "candidates"
}

func (c *Cand) Reset() {
	*c = Cand{}
}

func (c *Cand) SetActive(b bool) {
	c.Active = b
}

func (c Cand) GetSelf() interface{} {
	return c
}

func (c *Cand) GetId() int {
	return c.Id
}

func (c *Cand) SetId(id int) {
	c.Id = id
}

//Search only candidates with active true status and given q string and returns RCands array if
//Example: Cand{}.Search("Candidate Name | Department Name", 10)
func (c Cand) Search(q string, limit int) (r []RCands, err error) {

	sql := `SELECT t.id, t.name, t.lname, d.title, t.phone, t.email, t.addr, t.salary, t.currency
	FROM candidates t LEFT JOIN departments d ON t.ref_dep=d.id
	WHERE t.active=TRUE AND (to_tsvector(t.name)||to_tsvector(t.lname)||to_tsvector(d.title)) @@ to_tsquery($1)
	ORDER BY id LIMIT $2`

	_, err = Select(&r, sql, strings.Replace(strings.Trim(q, " "), " ", " & ", -1), limit)

	return
}
