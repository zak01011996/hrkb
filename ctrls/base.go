package ctrls

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	config "hrkb/conf"
	M "hrkb/models"
	"hrkb/services/bot"
	"hrkb/services/notify"
	"hrkb/utils"

	"github.com/astaxie/beego"
	"github.com/nicksnyder/go-i18n/i18n"
)

type CtrlErr string

var (
	T                        i18n.TranslateFunc
	DM                       *M.DM
	staticDir                string
	langTypes                map[string]string
	NFS                      *notify.NotificationService
	gitlabUrl, gitlabProject string
	notNilErr                error   = errors.New("there was error")
	internalErr              CtrlErr = "internal_error"
	dbErr                    CtrlErr = "db_error"
	marshalErr               CtrlErr = "marshal_error"
	gitlabAuthFailed         CtrlErr = "gitlab_auth_failed"
	gitlabTokenReq           CtrlErr = "gitlab_token_required"
	Bot                      *bot.Bot
)

//Struct for responce with json data
type RJson struct {
	Msg     string `json:"message"`
	Success bool   `json:"success"`
}

type BaseController struct {
	beego.Controller
	staticDir string
}

func CtrlInit() {

	staticDir = beego.AppConfig.String("staticDir")

	CandDateFmt = beego.AppConfig.String("comments::cand_date_format")
	TrashDateFmt = beego.AppConfig.String("comments::trash_date_format")
	gitlabUrl = beego.AppConfig.String("gitlab::url")
	gitlabProject = beego.AppConfig.String("gitlab::project")

	beego.AddFuncMap("urlFor", beego.UrlFor)
	beego.AddFuncMap("js", js)
	beego.AddFuncMap("css", css)
	beego.AddFuncMap("kb", kb)
	beego.AddFuncMap("isImg", isImg)
	beego.AddFuncMap("cutStr", CutLongText)
	beego.AddFuncMap("urlContain", urlContain)
	beego.AddFuncMap("T", i18n.IdentityTfunc)
	beego.AddFuncMap("floatStr", floatToString)

	DM = M.GetDM()

	// init notify service
	initNotifyService()
	initBot()

}
func (c *BaseController) CheckErr(err error, ctrlErr CtrlErr, msg ...string) bool {

	if err == nil {
		return false
	}

	c.Data["json"] = RJson{Msg: T(string(ctrlErr))}
	beego.Error(c.Ctx.Request.URL, ctrlErr, msg, err)
	return true
}

func (c *BaseController) SendMails(html string) {
	var users []M.User

	err := DM.FindAll(&M.User{}, &users, M.Sf{}, M.Where{And: M.W{"NotifyByMail": true}})

	if err != nil {
		beego.Error("Find users error ", err)
	}

	mail := M.Mail{}
	mail.ToType = "to"

	//Getting config data
	mail.FromMail = beego.AppConfig.String("mail::from_mail")
	mail.FromName = sql.NullString{beego.AppConfig.String("mail::from_name"), true}
	mail.Subject = sql.NullString{beego.AppConfig.String("mail::from_subject"), true}

	mail.Html = sql.NullString{html, true}
	mail.Text = sql.NullString{"", true}

	for _, user := range users {
		mail.ToName = sql.NullString{user.Name, true}
		mail.ToMail = user.Email
		//Send with Notification Service
		NFS.Send(mail)
	}
}

//Sends Telegram notifications to the group
func (c *BaseController) SendTgNotifs(s string) {
	NFS.Send(s)
}

func LoadLangs() {
	var langs []M.Lang
	l := &Lang{}
	err := DM.FindAll(&M.Lang{}, &langs, M.Sf{"Code"}, M.Where{}, M.NewParams(M.Params{Sort: "Code ASC"}))
	if err != nil {
		beego.Error(err)
		return
	}
	for _, v := range langs {
		i18n.MustLoadTranslationFile(l.langFileName("lang::folder", v.Code))
	}
}

func (c *BaseController) Prepare() {

	var lang string
	var langs []M.Lang
	var dLang string
	langTypes := make(map[string]string)

	c.Data["title"] = "Human Resource Management"
	c.Data["scripts"] = []string{}
	c.Data["styles"] = []string{}
	c.Data["url"] = c.Ctx.Request.RequestURI
	c.Layout = "layout.tpl"

	if c.IsAjax() {
		c.EnableRender = false
	}

	if uid := c.GetSession("uid"); uid != nil {
		c.Data["gitlabToken"] = c.GetSession("gitlabToken")
		c.Data["UserName"] = c.GetSession("name")
		c.Data["UserIsAdmin"] = c.GetSession("role").(int) == M.RoleAdmin
		c.Data["UserId"] = uid.(int)
	}

	flash := beego.ReadFromRequest(&c.Controller)

	if n, ok := flash.Data["notice"]; ok {
		c.Data["notice"] = n
	}

	c.Data["rUrl"] = c.Ctx.Request.RequestURI

	err := DM.FindAll(&M.Lang{}, &langs, M.Sf{"Code", "IsDefault"}, M.Where{}, M.NewParams(M.Params{Sort: "Code ASC"}))
	if err != nil {
		beego.Error(err)
	}
	for _, v := range langs {
		if v.IsDefault {
			dLang = v.Code
		}
		s := v.Code[:2]
		langTypes[v.Code] = s
	}

	uLang := c.Input().Get("lang")

	if uLang != "" && langTypes[uLang] != "" {
		c.SetSession("lang", uLang)
	} else if s := c.GetSession("lang"); s != nil && langTypes[s.(string)] != "" {
		uLang = s.(string)
	}

	aLang := c.Ctx.Request.Header.Get("Accept-Language")

	if len(aLang) > 4 {
		aLang = strings.ToLower(aLang[:3]) + strings.ToUpper(aLang[3:5])
		if _, ok := langTypes[aLang]; !ok {
			aLang = ""
		}
	} else {
		aLang = ""
	}

	Tfn, err := i18n.Tfunc(uLang, aLang, dLang)
	T = Tfn
	M.T = Tfn
	if err != nil {
		beego.Error(err)
	}

	beego.AddFuncMap("T", Tfn)

	switch {
	case langTypes[uLang] != "":
		lang = uLang
	case langTypes[aLang] != "":
		lang = aLang
	default:
		lang = dLang
	}

	c.Data["Langs"] = langTypes
	c.Data["CurrentLang"] = lang
}

func (c *BaseController) AddAssets(index string, assets []string) {
	ss := c.Data[index].([]string)

	for _, asset := range assets {
		if i := utils.IndexOfStr(ss, asset); i == -1 {
			ss = append(ss, asset)
		}
	}
	c.Data[index] = ss
}

func (c *BaseController) Finish() {
	if c.IsAjax() {
		c.ServeJson()
	}
}

// Initiate nofication service
func initNotifyService() {
	//Init Notification Service Config
	var (
		err     error
		groupId int
	)

	groupId, err = beego.AppConfig.Int("telegram::group_id")

	if err != nil {
		beego.Error("Telegram config read error: group_id", err)
	}

	conf := notify.NotifyMailConfig{}

	err = config.ParseConfig(beego.AppConfig, &conf)
	if err != nil {
		log.Fatalf("Mail Config Read Error %s", err)
		return
	}

	conf.DB = DM

	//Init Notification Service
	NFS = notify.NewNotificationService(notify.Conf{conf, conf.BuffLimit, groupId})
	NFS.Start()
	//Starting Notification Service logging
	go ChannLogger(NFS.Errors(), NFS.Notifications())

}

//Channel Logger for listening service errors and info and log them
func ChannLogger(errChan <-chan error, infoChan <-chan string) {
	for {
		select {
		case err := <-errChan:
			beego.Error(err)
		case info := <-infoChan:
			beego.Info(info)
		}
	}
}

// !!! TELEGRAM SERVICES NEED CONFIGURED TELEGRAM-CLI
// 1.Download and compile telegram cli on your machine https://github.com/vysheng/tg
// 2.Authentificate manually
// 3.Put your public-key to /etc/telegram-cli/server.pub
//
//Initializes HrkbBot for listening and responding commands
func initBot() {
	botConf := bot.Conf{}
	err := config.ParseConfig(beego.AppConfig, &botConf)

	if err != nil {
		log.Fatalf("Bot Config Read Error %s", err)
		return
	}

	botConf.Cmd = NFS.TgService.GetCmd()

	Bot = bot.NewBot(botConf)
	Bot.Start()

	go ChannLogger(Bot.Errors(), Bot.Notifications())
}
