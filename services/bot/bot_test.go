package bot

import (
	"hrkb/services/notify"
	"log"
	"testing"
	"time"

	beeconf "github.com/astaxie/beego/config"
	"github.com/nicksnyder/go-i18n/i18n"
	"hrkb/conf"
	M "hrkb/models"
)

var (
	botConf Conf
)

func init() {

	cmd, err := notify.NewCmd("telegram-cli", "-C")

	if err != nil {
		log.Fatal("NewCmd error", err)
	}

	botConf = Conf{
		Limit: 10,
		Url:   "localhost",
		Rpms:  500,
		Cmd:   cmd,
	}

	//Init Data Mapper
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

	M.PrepareTables(&M.Cand{})
	Tfn, _ := i18n.Tfunc("en-us", "en-us", "en-us")
	M.T = Tfn
}

type TestReader struct{}

func (t *TestReader) Read(p []byte) (n int, err error) {
	str := "[15:42]  GroupName User Name >>> #find test \n [15:42]  GroupName User Name >>> #find abdullo \n [15:50]  Abdullo Xidoyatov >>> #salom\n"
	n = copy(p[:], str)
	return
}

func (t *TestReader) Close() error {
	return nil
}

func TestStart(t *testing.T) {
	reader := &TestReader{}

	bot := NewBot(botConf)
	bot.Cmd.Reader = reader

	bot.Start()

	select {
	case <-bot.Notifications():
	case <-bot.Errors():
		t.Error("Bot start error")
	case <-time.After(time.Second):
		t.Error("Bot start timeout")
	}
}
