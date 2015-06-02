package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/nicksnyder/go-i18n/i18n"
	"hrkb/conf"
	"hrkb/utils"
)

type Model interface {
	Table() string
	Reset()
	SetActive(b bool)
	GetSelf() interface{}
	GetId() int
	SetId(id int)
}

type W map[string]interface{}

type Where struct {
	And, Or W
}

type Sf []string

type pvalues []interface{}

type Params struct {
	Offset, Limit int
	Sort          string
}

type Transaction struct {
	t *gorp.Transaction
}

var (
	dbmap *gorp.DbMap
	dm    *DM
	dic   dictype
	T     i18n.TranslateFunc
)

type sdictype map[string]string
type dictype map[string]sdictype

type ValidMap struct {
	validation.Validation
}

type ValidMapFormer interface {
	Valid(*ValidMap)
}

func transErrors() {
	for k, _ := range validation.MessageTmpls {
		validation.MessageTmpls[k] = T("valid_" + k)
	}
}

func DbOpen(dbConf conf.Database) error {

	var dial gorp.Dialect
	var dsn string

	switch dbConf.Driver {
	case "postgres":
		dial = &gorp.PostgresDialect{}
		dsn = "host=" + dbConf.Postgres.Host + " dbname=" + dbConf.Postgres.Database + " user=" + dbConf.Postgres.User + " password=" + dbConf.Postgres.Pass
	default:
		return errors.New(fmt.Sprintf("unknown driver \"%s\"", dbConf.Driver))
	}

	db, err := sql.Open(dbConf.Driver, dsn)

	if err != nil {
		return err
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: dial}

	dm = NewDM(dbmap)

	return nil
}

func DbClose() {
	dbmap.Db.Close()
}

func GetDM() *DM {
	return dm
}

/*
  Select from expressin field and operator
  example param1="Login>" return "Login",">"
  example param1="Age<18" return "Age","<"
  example param1="Id" return "Id","="
*/

func splitOp(f string) (s string, op string) {

	op = "="
	s = f

	switch {
	case strings.Contains(f, ">"):
		op = ">"
		s = s[:len(s)-1]
	case strings.Contains(f, "<"):
		op = "<"
		s = s[:len(s)-1]
	case strings.Contains(f, "!="):
		op = "!="
		s = s[:len(s)-2]
	}

	return
}

func prepActiveCond(w Where) (a W, o W) {

	a, o = w.And, w.Or

	if _, ok := a["Active"]; ok {
		return
	}

	if _, ok := o["Active"]; ok {
		return
	}

	if a == nil {
		a = make(W, 0)
	}

	a["Active"] = true

	if o == nil || len(o) == 0 {
		return
	}

	o["Active"] = true

	return

}

/*
generate SQL query
param1 Model interface
param2 Sf Fields to select empty=*
param3 AndWhere conditions key1=value1,key2=value2
param4 OrWhere conditions key1=value1,key2=value2
param5 optional additional query parameters (sort, limit etc.). Can be generated NewParams function

	prepSelect(&User{},Sf{"Id","Login"},W{"Id":1},W{},NewParams(Params{Limit:5}))
        SELECT id,login FROM users WHERE id=$1 LIMIT $2
        [1,5]
*/

func prepSelect(table string, dic sdictype, sf Sf, andW, orW W, p *Params) (sql string, pv pvalues) {

	var s, w string

	for i, j := 0, len(sf); i < j; i += 1 {
		if val, ok := dic[sf[i]]; ok {
			s += val + ","
		}
	}

	if s != "" {
		s = s[:len(s)-1]
	} else {
		s = "*"
	}

	w = concatConds("OR", getCmpCond(andW, &pv, dic, "AND"), getCmpCond(orW, &pv, dic, "AND"))

	sql = "SELECT " + s + " FROM " + table
	if w != "" {
		sql += " WHERE " + w
	}

	if p != nil {
		sql += p.buildSql(dic, &pv)
	}

	return
}

/*
 find single row from db
 params model,sf,andw,orw,p equals buildSql function
 return bool if row is founded in db
*/
func find(m Model, sf Sf, andW, orW W, p ...*Params) error {

	var p2 *Params
	if len(p) > 0 {
		p2 = p[0]
		p2.Limit = 1
	} else {
		p2 = &Params{Limit: 1}
	}

	table := m.Table()

	sql, pv := prepSelect(table, GetTableDic(table), sf, andW, orW, p2)

	return dbmap.SelectOne(m, sql, pv...)

}

/*
 find multi rows from db
 params model,sf,andw,orw,p equals buildSql function
 param h is container for result
 return true if h is not empty
*/
func findAll(m Model, h interface{}, sf Sf, andW, orW W, p ...*Params) error {

	var p2 *Params
	if len(p) > 0 {
		p2 = p[0]
	} else {
		p2 = NewParams(Params{})
	}

	table := m.Table()

	sql, pv := prepSelect(table, GetTableDic(table), sf, andW, orW, p2)

	_, err := dbmap.Select(h, sql, pv...)

	return err

}

func NewParams(p Params) *Params {
	var np Params = p

	if np.Sort == "" {
		np.Sort = "Id ASC"
	}

	return &np
}

/*
build extra sql from based Params structure
param1 Model
param2 pvalues
		object=Params{Limit:5}
		param1=&User{}
		param2=[10,20]
		return " LIMIT $3"
		and pvalues modified to [10,20,5]

*/
func (p *Params) buildSql(dic sdictype, pv *pvalues) (sql string) {

	if p.Sort != "" {
		sort := p.Sort
		for k, v := range dic {
			sort = strings.Replace(sort, k, v, -1)
		}
		sql = " ORDER BY " + sort
	}

	if p.Limit > 0 {
		sql += " LIMIT $" + strconv.Itoa(len(*pv)+1)
		*pv = append(*pv, p.Limit)
	}

	if p.Offset > 0 {
		sql += " OFFSET $" + strconv.Itoa(len(*pv)+1)
		*pv = append(*pv, p.Offset)
	}

	return
}

/**
	Execute INSERT SQL Query and fill model from execution result
	m: model
	f: inserted fields (optional)
		insert(&User{Login:"abc",Password:"123",Role:"XXX"},Sf{"Login"}) result:  INSERT INTO users (login) VALUES ("abc")
**/

func insert(m Model, f Sf) error {

	var pv pvalues

	if len(f) == 0 {
		return dbmap.Insert(m)
	}

	fs, fv, n, table := "", "", 0, m.Table()

	d := GetTableDic(table)

	r := reflect.ValueOf(m.GetSelf())

	for _, v := range f {

		if val, ok := d[v]; ok {

			f := r.FieldByName(v)
			if f.IsValid() {
				n += 1
				pv = append(pv, f.Interface())
				fs += val + ", "
				fv += "$" + strconv.Itoa(n) + ", "
			}

		}
	}

	if fs == "" {
		return errors.New("no fields to insert")
	}

	sql := "INSERT INTO " + table + " (" + utils.CutStr(fs, 2) + ") VALUES (" + utils.CutStr(fv, 2) + ") RETURNING *"
	return dbmap.SelectOne(m, sql, pv...)
}

/**
	generate and execute UPDATE SQL Query
	m model
	f updating fields
	returned affected_rows_count,error (if occured)
		update(&User{Login:"new_login",Id:10},Sf{"Login","Surname"}) result: UPDATE users SET login=new_login WHERE id=10
**/
func update(m Model, f Sf) (int64, error) {

	var pv pvalues

	if len(f) == 0 {
		return dbmap.Update(m)
	}

	fs, n, table := "", 0, m.Table()

	d := GetTableDic(table)

	r := reflect.ValueOf(m.GetSelf()) //collect cloned object values

	for _, v := range f {

		if val, ok := d[v]; ok {

			f := r.FieldByName(v)
			if f.IsValid() {
				n += 1
				pv = append(pv, f.Interface())
				fs += val + "=$" + strconv.Itoa(n) + ", " //include object value to query if it present in f (updating fields)
			}

		}
	}

	if fs == "" {
		return 0, errors.New("no fields to update")
	}

	if val, ok := d["Id"]; ok {
		n += 1
		pv = append(pv, m.GetId())
		sql := "UPDATE " + table + " SET " + utils.CutStr(fs, 2) + " WHERE " + val + "=$" + strconv.Itoa(n)

		r, err := dbmap.Exec(sql, pv...)

		if err != nil {
			return 0, err
		}

		return r.RowsAffected()
	}

	return 0, errors.New("id not found")

}

func exec(query string, args ...interface{}) (sql.Result, error) {
	return dbmap.Exec(query, args...)
}

func Begin() (*Transaction, error) {
	t, err := dbmap.Begin()
	return &Transaction{t: t}, err
}

func (t *Transaction) Commit() error {
	return t.t.Commit()
}

func (t *Transaction) Rollback() error {
	return t.t.Rollback()
}

func (t *Transaction) Insert(list ...interface{}) error {
	return t.t.Insert(list...)
}

func (t *Transaction) Update(list ...interface{}) (int64, error) {
	return t.t.Update(list...)
}

func (t *Transaction) Delete(list ...interface{}) (int64, error) {
	return t.t.Delete(list...)
}

func (t *Transaction) exec(query string, args ...interface{}) (sql.Result, error) {
	return t.t.Exec(query, args...)
}

/*
	run some actions for table models
	called 1 time on application starting and before running tests
*/
func PrepareTables(args ...Model) {
	for _, m := range args {
		dbmap.AddTableWithName(m.GetSelf(), m.Table()).SetKeys(true, "id")
		GenDic(m)
	}
}

/*
  Generate dictonary for models (param1)
  generated dictonary exist in 1 exemplar per model
  dictonary contain map when key=struct property, and value = db column name
  generated dictonary can be called as dic[<table>][<property>]
  example
  		GenDic(&Users{},&Crit{})
  		dic["users"]["Login"] //return login
  		dic["criteria"]["Cat"] //return ref_dep
*/
func GenDic(m ...Model) {

	if dic == nil {
		dic = make(dictype, 0) //if GenDic is called in first time, initialize variable
	}

	for i, j := 0, len(m); i < j; i += 1 {

		table := m[i].Table()
		if _, ok := dic[table]; !ok { //if not generated dictonary for this model generate it
			dic[table] = make(sdictype, 0)
			s := reflect.ValueOf(m[i]).Elem()
			typeOfT := s.Type()
			for n, m := 0, s.NumField(); n < m; n += 1 { //iterate all properties from model struct
				f := typeOfT.Field(n)
				k := strings.Split(f.Tag.Get("db"), ";")
				if len(k) > 0 && k[0] != "-" {
					dic[table][f.Name] = k[0] //first value from db tag saves to dictonary as column name
				}
			}
		}
	}
}

/*
	Model validation
	operated with model property "valid" tag and function "Valid" if it defined in model
	param1 = pointer to model
	retured ValidMap type result

	example
	  u := User{Password:"123",PassConf:"456"}
	  v := Validate(&u)

	  v.HasErrors() //true
	  v.Errors // contain PassConf = Passwords do not match
*/
func Validate(m interface{}) ValidMap {

	transErrors()

	v := ValidMap{}

	v.Valid(m) //call beego Validation mechanism ;)

	if !v.HasErrors() {
		if form, ok := m.(ValidMapFormer); ok {
			form.Valid(&v)
		}
	}

	return v
}

func GetTableDic(table string) (d sdictype) {
	if _, ok := dic[table]; ok {
		d = dic[table]
	}
	return
}

func ExpandFormErrors(v *ValidMap, data map[interface{}]interface{}) {

	if !v.HasErrors() {
		return
	}

	for _, err := range v.Errors {
		data["err"+strings.Split(err.Key, ".")[0]] = err.Message
	}
}

func Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return dbmap.Select(i, query, args...)
}

func count(m Model, andW, orW W) (int64, error) {

	var pv pvalues

	table := m.Table()
	dic := GetTableDic(table)

	w := concatConds("OR", getCmpCond(andW, &pv, dic, "AND"), getCmpCond(orW, &pv, dic, "AND"))

	sql := "SELECT COUNT(*) FROM " + table
	if w != "" {
		sql += " WHERE " + w
	}

	return dbmap.SelectInt(sql, pv...)
}

/**
 build condition sql string by conditions
 		l: conditions
 		pv: container for excluded values
 		dic: table dictonary
 		cmp_op: compare operator between conditions

 			getCmpCond(W{"Id":10,"Login":"admin"},pv,users_dic,"AND") returns (id=$1 AND login=$2), and stores to pv
 			10,admin
**/
func getCmpCond(l W, pv *pvalues, dic sdictype, cmp_op string) (sql string) {

	var n int = len(*pv)

	for k, v := range l {
		k, op := splitOp(k)
		if val, ok := dic[k]; ok {
			n += 1
			sql += val + op + "$" + strconv.Itoa(n) + " " + cmp_op + " "
			*pv = append(*pv, v)
		}
	}

	if i := len(sql); i > 0 {
		sql = "(" + sql[:i-len(cmp_op)-2] + ")"
	}

	return
}

/**
 concate strings with defined operator between
 		concatCOnds("OR","a","b","c","d") returns a OR b OR c OR d
**/
func concatConds(op string, args ...string) (s string) {

	for _, v := range args {
		if v != "" {
			s += v + " " + op + " "
		}
	}

	if s != "" {
		s = utils.CutStr(s, len(op)+2)
	}

	return
}
