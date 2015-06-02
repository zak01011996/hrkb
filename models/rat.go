package models

type Rat struct {
	Id     int  `db:"id" form:"id"`
	Crit   int  `db:"ref_crit" valid:"Required" form:"crit"`
	Cand   int  `db:"ref_cand" valid:"Required" form:"cand"`
	User   int  `db:"ref_user"`
	Value  int  `db:"value" valid:"Required;Min(1);Max(5)" form:"value"`
	Active bool `db:"active"`
}

type GroupedRating struct {
	Id       int     `db:"id"`
	Criteria string  `db:"criteria"`
	Rating   float64 `db:"rating"`
	Dep      string  `db:"dep"`
}

type DetailedRating struct {
	GroupedRating
	Author string `db:"author"`
}

type UserRating struct {
	DetailedRating
	AuthorId int `db:"author_id"`
}

func (u Rat) Table() string {
	return "ratings"
}

func (r *Rat) Reset() {
	*r = Rat{}
}

func (r *Rat) SetActive(b bool) {
	r.Active = b
}

func (r Rat) GetSelf() interface{} {
	return r
}

func (r *Rat) GetId() int {
	return r.Id
}

func (r *Rat) SetId(id int) {
	r.Id = id
}

func (r *Rat) GetAverageRating(c_id int) (float64, error) {
	v, err := dbmap.SelectFloat(
		"SELECT round(CAST(sum(value) AS DEC(5,2))/count(id), 2) as total FROM ratings WHERE ref_cand=$1",
		c_id)
	return v, err
}

func (r *Rat) GetGroupedRatings(c_id int) ([]GroupedRating, error) {

	var v []GroupedRating
	_, err := Select(&v,
		`SELECT c.id id, c.title criteria, d.title dep, round(CAST(sum(r.value) AS DEC(5,2))/count(r.id), 2) rating 
		 FROM ratings r 
		 LEFT JOIN criteria c ON c.id = r.ref_crit
		 LEFT JOIN departments d ON d.id = c.ref_dep
		 WHERE r.active = TRUE AND r.ref_cand = $1
		 GROUP BY c.id, c.title, d.title
		 ORDER BY c.title;`,
		c_id)
	return v, err
}

func (r *Rat) GetDetailedRatings(c_id int) ([]DetailedRating, error) {

	var v []DetailedRating
	_, err := Select(&v,
		`SELECT r.id id, c.title criteria, d.title dep, r.value rating, u.name author
		 FROM ratings r
		 LEFT JOIN users u ON u.id = r.ref_user
		 LEFT JOIN criteria c ON c.id = r.ref_crit
		 LEFT JOIN departments d ON d.id = c.ref_dep
		 WHERE r.active = TRUE AND r.ref_cand = $1
		 ORDER BY c.title;`,
		c_id)
	return v, err
}

func (r *Rat) GetUserRatings(c_id int, u_id int) ([]UserRating, error) {

	var v []UserRating
	_, err := Select(&v,
		`SELECT r.id id, c.title criteria, r.value rating, d.title dep, u.name author, u.id author_id
		 FROM ratings r
		 LEFT JOIN users u ON u.id = r.ref_user
		 LEFT JOIN criteria c ON c.id = r.ref_crit
		 LEFT JOIN departments d ON d.id = c.ref_dep
		 WHERE r.active = TRUE AND r.ref_cand = $1 AND u.id = $2
	         ORDER BY c.title;`,
		c_id, u_id)
	return v, err
}
