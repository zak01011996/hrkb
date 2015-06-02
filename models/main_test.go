package models

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/astaxie/beego/config"
	"github.com/nicksnyder/go-i18n/i18n"
	"hrkb/conf"
)

func InsertTestData(f string) error {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	_, err = dbmap.Exec(string(dat))
	return err
}

func ClearTestData(m ...Model) error {
	for _, v := range m {
		_, err := dbmap.Exec("TRUNCATE TABLE " + v.Table() + " CASCADE;")
		if err != nil {
			return err
		}

		_, err = dbmap.Exec("ALTER SEQUENCE " + v.Table() + "_id_seq RESTART WITH 1;")
		if err != nil {
			return err
		}

	}

	return nil
}

func TestMain(m *testing.M) {

	c, err := config.NewConfig("ini", "../conf/app.conf")

	if err != nil {
		log.Fatal(err)
	}

	TestConfig, err := conf.Initialize("test", c)
	if err != nil {
		log.Fatal(err)
	}

	if err := DbOpen(TestConfig.Db); err != nil {
		log.Fatal(err)
	}

	defer DbClose()

	PrepareTables(&User{}, &Cand{}, &Rat{}, &Crit{}, &Dep{}, &Comment{}, &Role{}, &Contact{}, &Mail{})

	err = ClearTestData(&User{}, &Crit{}, &Rat{}, &Dep{})
	if err != nil {
		log.Fatal(err)
	}

	err = InsertTestData("../migrations/test_data/data.sql")
	if err != nil {
		log.Fatal(err)
	}

	i18n.MustLoadTranslationFile("../langs/en-US.all.json")
	T, _ = i18n.Tfunc("en-US", "en-US", "en-US")

	exitCode := m.Run()

	os.Exit(exitCode)
}
