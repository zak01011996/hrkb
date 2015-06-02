package bot

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	M "hrkb/models"
	"hrkb/services/notify"
	"regexp"
	"strings"
)

const (
	IN               string = "<<<"
	OUT              string = ">>>"
	CMDSIGN          string = "#"
	BOTCMD_NOT_FOUND string = "Bot command not found"
	TG_EXTRA_HEXCODE string = "\x1b[K>"
	BOTNAME          string = "HRBot"
)

var (
	//Base Url of application get by by Bot Config.
	//It used for generate link to candidate page
	Url string
)

//Bot configuration struct
type Conf struct {
	Limit int    `conf:"bot::buff_limit"`
	Url   string `conf:"bot::url"`
	Rpms  int    `conf:"bot::rate_per_msec"`
	Cmd   *notify.Cmd
}

//Bot commands
type BotCmd struct {
	From   string   //Telegram: GroupName | Username
	User   string   //Telegram: User Name
	Where  string   //Telegram: >>> | <<<
	Params []string // Telegram: message body | params
}

//Every bot function recieves BotCmd and responds with string array
type BotFunc func(c BotCmd, errChan notify.ErrorChannel) []string

//Bot functions
//Greets with sender
func cmdSalom(c BotCmd, errChan notify.ErrorChannel) []string {
	return []string{fmt.Sprintf("Salom %s", strings.Split(c.User, "_")[0])}
}

//Gets params from command and send it to Cand search
func cmdFind(c BotCmd, errChan notify.ErrorChannel) []string {
	//Send Params from 1:~ (params[0]  is command "#find")

	cands, err := M.Cand{}.Search(strings.Join(c.Params[1:], " "), 5)
	//Error while searching candidate

	if err != nil {
		errChan.In() <- errors.New(fmt.Sprintln("Candidate search error ", err))
	}

	str := []string{}

	//If nothing found
	if len(cands) < 1 {
		str = append(str, fmt.Sprintf("Sorry I cant find anything for query \"%s\" :-(. But I bealive he/she will be added", strings.Join(c.Params[1:], " ")))
		return str
	}

	//If found only one candidate return detailed data about cand
	if len(cands) == 1 {
		cand := cands[0]
		str = append(str, fmt.Sprintf("Name: %s %s: ", cand.Name, cand.LName))
		str = append(str, fmt.Sprintf("Tel: %s", cand.Tel))
		str = append(str, fmt.Sprintf("Email: %s", cand.Mail))
		str = append(str, fmt.Sprintf("Department: %s", cand.Dep))
		str = append(str, fmt.Sprintf("Address: %s", cand.Addr))
		str = append(str, fmt.Sprintf("Salary: %s %s", cand.Salary, cand.Currency))
		return str
	}

	//Loop through found candidates and return them
	for i, cand := range cands {
		str = append(str, fmt.Sprintf("%d. %s %s. Tel: %s Email: %s", i+1, cand.Name, cand.LName, cand.Tel, cand.Mail))
	}
	return str
}

type Bot struct {
	Cmd       *notify.Cmd                //telegram-cli command wrappper
	closeChan chan bool                  //close chan for quitting from bot
	cmdRegexp *regexp.Regexp             //Regular expression for parsing telegram messages
	rpms      int                        //Rate per Milli second for listen commands
	Commands  map[string]BotFunc         //Commands for bot
	ntfChan   notify.NotificationChannel //Notification Channel for sending notificatation outside
	errChan   notify.ErrorChannel        //For Sending errors outside
	url       string                     //For sending candidate urls
}

func NewBot(c Conf) *Bot {
	// Regular Expression matche format
	// 	[15:42]  ChatName User Name IN CMDSIGN message body params
	//	!!! IN & CMDSIGN are constants
	//
	// Ex(default):
	// 	[15:42]  HRKB Abdullo Xidoyatov >>> #test second
	// 	[15:50]  Abdullo Xidoyatov >>> test second third fourth
	// 	[15:55]  Abdullo Xidoyatov <<< test second third
	// 	[15:56]  HRKB HRBot >>> test second third
	//
	// And it have submatch groups
	//	[1] from
	//	[2] where
	// 	[3] body

	regStr := fmt.Sprintf(`\[\d{2}:\d{2}\](?P<from>(?:\s+\w*){1,3}\s)(?P<where>%s|%s)(?P<body>(?:\s%s\w*|\s[a-zA-ZА-Яа-я]*)+)`, IN, OUT, CMDSIGN)

	regExp := regexp.MustCompile(regStr)

	cmds := make(map[string]BotFunc)

	cmds[CMDSIGN+"salom"] = cmdSalom
	cmds[CMDSIGN+"find"] = cmdFind

	Url = c.Url

	bot := &Bot{
		Cmd:       c.Cmd,
		closeChan: make(chan bool),
		cmdRegexp: regExp,
		Commands:  cmds,
		ntfChan:   notify.NewNotificationChannel(c.Limit),
		errChan:   notify.NewErrorChannel(c.Limit),
		rpms:      c.Rpms,
		url:       c.Url,
	}

	return bot
}

func (b *Bot) ParseCommand(out string, cmds *[]BotCmd) {

	res := b.cmdRegexp.FindAllStringSubmatch(out, -1)

	//If no commands parsed return not found error
	if len(res) <= 0 {
		return
	}

	for _, found := range res {

		botc := BotCmd{}
		//Getting from group after parsing command
		from := strings.Split(strings.TrimSpace(found[1]), " ")

		if len(from) < 3 {
			//If len less than 3 this message from user
			//Ex:
			//	from["User", "Name"]
			//	from["UserName"]
			botc.From = strings.Join(from, "_")
			botc.User = strings.Join(from, "_")
		} else {
			//Else this message from group
			//Ex: from["GroupName", "User", "Name"]
			botc.User = fmt.Sprintf("%s_%s", from[1], from[2])
			botc.From = fmt.Sprintf("%s", from[0])
		}

		//Getting where group
		where := strings.TrimSpace(found[2])
		botc.Where = where

		//Getting body of the message
		body := strings.Split(strings.TrimSpace(found[3]), " ")
		botc.Params = body

		*cmds = append(*cmds, botc)
	}
}

//Remove redundant needle from output
func (b *Bot) rmExtraHexcode(str, needle string) string {
	return strings.Replace(str, needle, "", -1)
}

//Start bot and listen in n time tick for commands
func (b *Bot) Start() {
	scanner := bufio.NewScanner(b.Cmd.Reader)
	go func() {
		//Scan Telgram putput
		for scanner.Scan() {
			//Reading runes
			runes := bytes.Runes(scanner.Bytes())
			//Convert to string
			out := string(runes[:])
			//Replace extra hexcode
			output := b.rmExtraHexcode(out, TG_EXTRA_HEXCODE)

			cmds := []BotCmd{}
			//Parse Commands
			b.ParseCommand(output, &cmds)

			//If commands parsed succesfully
			//Loop throug all parsed commands
			if len(cmds) == 0 {
				continue
			}

			//Loop throug all parsed commands
			for _, cmd := range cmds {
				//If command found in Commands list send its return
				if fn, ok := b.Commands[cmd.Params[0]]; ok {
					//Send multi line response
					for _, r := range fn(cmd, b.errChan) {
						b.ntfChan.In() <- fmt.Sprintf("Cmd: %s", cmd)
						b.Cmd.Input(fmt.Sprintf("msg %s %s", cmd.From, r))
					}
				}
			}
		}

	}()
}

func (b *Bot) Notifications() <-chan string {
	return b.ntfChan.Out()
}

func (b *Bot) Errors() <-chan error {
	return b.errChan.Out()
}
