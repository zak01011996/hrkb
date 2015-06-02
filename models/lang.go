package models

type (
	Lang struct {
		Id        int    `db:"id"`
		Code      string `db:"code" form:"code" valid:"Required"`
		IsDefault bool   `db:"is_default"`
		Active    bool   `db:"active"`
	}
)

var (
	LangsSupports = map[string]string{
		"ar-AR": "Arabic",
		"be-BE": "Belarusian",
		"bg-BG": "Bulgarian",
		"ca-CA": "Catalan",
		"zh-ZH": "Chinese",
		"cs-CS": "Czech",
		"da-DA": "Danish",
		"nl-NL": "Dutch",
		"en-US": "English",
		"fr-FR": "French",
		"de-DE": "German",
		"is-IS": "Icelandic",
		"it-IT": "Italian",
		"ko-KO": "Korean",
		"ja-JA": "Japanese",
		"lt-LT": "Lithuanian",
		"pl-PL": "Polish",
		"pt-PT": "Portuguese",
		"pt-BR": "Portuguese (Brazilian)",
		"ru-RU": "Russian",
		"es-ES": "Spanish",
		"sv-SV": "Swedish",
		"uk-UK": "Ukrainian",
		"uz-UZ": "Uzbekistan",
	}
)

func (l Lang) Table() string {
	return "langs"
}

func (l *Lang) Reset() {
	*l = Lang{}
}

func (l *Lang) SetActive(b bool) {
	l.Active = b
}

func (l Lang) GetSelf() interface{} {
	return l
}

func (l *Lang) GetId() int {
	return l.Id
}

func (l *Lang) SetId(id int) {
	l.Id = id
}

func (l Lang) SetDefault(id int) error {
	_, err := exec("UPDATE langs SET is_default=$1 WHERE is_default=$2", false, true)
	if err == nil {
		_, err = update(&Lang{IsDefault: true, Id: id}, Sf{"IsDefault"})
	}
	return err
}

func (l *Lang) Valid(v *ValidMap) {
	if _, ok := LangsSupports[l.Code]; !ok {
		v.SetError("Code", T("lang_undefinedcode"))
	}
}
