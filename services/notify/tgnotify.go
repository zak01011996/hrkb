package notify

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type Cmd struct {
	CMD    *exec.Cmd
	Out    bytes.Buffer
	In     io.WriteCloser
	Reader io.ReadCloser
}

//Write string to input in bytes
func (c *Cmd) Input(s string) {
	io.WriteString(c.In, fmt.Sprintf("%s \n", s))
}

func NewCmd(cmdname string, opts ...string) (*Cmd, error) {

	//Initializing Command
	cmd := exec.Command(cmdname, opts...)

	//Creating Pipe for writing after iniztialization
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return &Cmd{}, err
	}

	//Creating Pipe for reading
	outPipe, err := cmd.StdoutPipe()

	if err != nil {
		return &Cmd{}, err
	}

	//Iniztialing command object
	c := &Cmd{
		CMD:    cmd,
		Out:    bytes.Buffer{},
		In:     inPipe,
		Reader: outPipe,
	}

	return c, nil
}

type NotifyTelegram struct {
	sendChan  chan string
	closeChan chan bool
	ntfChan   NotificationChannel
	errChan   ErrorChannel
	Cmd       *Cmd
	To        int
}

func NewNotifyTelegram(limit int, to int, nChan NotificationChannel, eChan ErrorChannel) *NotifyTelegram {
	cmd, err := NewCmd("telegram-cli", "-C")

	if err != nil {
		eChan <- err
	}

	return &NotifyTelegram{
		sendChan:  make(chan string, limit),
		closeChan: make(chan bool),
		ntfChan:   nChan,
		errChan:   eChan,
		Cmd:       cmd,
		To:        to,
	}
}

func (n *NotifyTelegram) Start() {
	n.ntfChan.In() <- "Telegram Service Started"

	go func() {
		n.Cmd.Input("dialog_list")
		if err := n.Cmd.CMD.Start(); err != nil {
			n.errChan.In() <- err
		}
		// add ping to tg-cli to keep alive connect
		ping := time.NewTicker(30 * time.Second)
		for {
			select {
			case ntf := <-n.sendChan:
				n.Cmd.Input(fmt.Sprintf("msg chat#%d %s", n.To, ntf))
				n.ntfChan.In() <- fmt.Sprintf("TgNotif send Success to: %d", n.To)
			case <-ping.C:
				n.Cmd.Input("dialog_list")
				n.ntfChan.In() <- "TgNoTif ping with dialog_list cmd"
			case <-n.closeChan:
				return
			}
		}
	}()
}

func (n *NotifyTelegram) Notify(s string) {
	n.sendChan <- s
}

func (n *NotifyTelegram) GetCmd() *Cmd {
	return n.Cmd
}

func (n *NotifyTelegram) Stop() {
	n.closeChan <- true
}
