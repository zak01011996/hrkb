package notify

import (
	beeconf "github.com/astaxie/beego/config"
	"github.com/nicksnyder/go-i18n/i18n"
	"hrkb/conf"
	M "hrkb/models"

	"fmt"
	"log"
	"testing"
)

var (
	db     *M.DM
	mailSt MailStore
)

func init() {
	c, err := beeconf.NewConfig("ini", "../../conf/app.conf")

	if err != nil {
		log.Fatal(err)
	}

	TestConfig, err := conf.Initialize("test", c)
	if err != nil {
		log.Fatal(err)
	}

	if err := M.DbOpen(TestConfig.Db); err != nil {
		log.Fatal(err)
	}

	M.PrepareTables(&M.Mail{})
	db = M.GetDM()
	Tfn, _ := i18n.Tfunc("en-us", "en-us", "en-us")
	M.T = Tfn

	mailSt = MailStore{db, 3}
}

func TestGetMail(t *testing.T) {

	mail := newMail(800)

	if err := mailSt.Insert(&mail); err != nil {
		t.Error("Insert error: ", err)
	}

	mail.Status = false
	mail.Try = 1

	if err := mailSt.Update(&mail); err != nil {
		t.Error("Save error: ", err)
	}

	mails := &[]M.Mail{}
	err := mailSt.GetFailed(mails, 100)

	if err != nil {
		t.Error("Get Mails error", err)
	}
}

func TestGetMailFail(t *testing.T) {

	mails := &[]M.Mail{}
	err := mailSt.GetFailed(mails, -1)

	if fmt.Sprintf("%s", err) != ErrLimit {
		t.Error("Limit error expected: ", err)
	}
}
