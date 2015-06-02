package notify

import (
	"errors"
	"fmt"
	"github.com/keighl/mandrill"
	M "hrkb/models"
	"strings"
	"time"
)

//Mail service instance
type NotifyMail struct {
	client    *mandrill.Client
	store     MailStore
	notifChan NotificationChannel
	sendChan  chan M.Mail
	errChan   ErrorChannel
	waitTime  int
	retry     int
	closeChan chan bool
	rpms      int
	frpm      int
}

//Mail Notifier config
type NotifyMailConfig struct {
	Apikey            string `conf:"mail::apikey"` //Mandrill Api Key
	DB                *M.DM  // Mail Store for manipulating with DB
	WaitTime          int    `conf:"mail::wait_time"`           // Wait time if internet connnection epsont wait for this time
	Timeout           int    `conf:"mail::timeout"`             // If Mandrill Api dont responce too long close connection in this time
	BuffLimit         int    `conf:"mail::buff_limit"`          // Buffer Limit
	Retry             int    `conf:"mail::retry"`               //If error occures while sending retry it n times
	RatePerH          int    `conf:"mail::hour_limit"`          //Mandrill Send limit per Hour
	RetryFailedPerMin int    `conf:"mail::failed_rate_per_min"` //Retry Sending failed messsages per n minutes
}

//Recieves Mail Service config
//And initializes ne NotifyMail obj and returns it
func NewNotifyMail(conf NotifyMailConfig, nChan NotificationChannel, eChan ErrorChannel) *NotifyMail {
	c := mandrill.ClientWithKey(conf.Apikey)

	c.HTTPClient.Timeout = time.Duration(conf.Timeout) * time.Second

	rps := 1

	//Computing rate Per Millisecond
	if conf.RatePerH > 3600000 {
		rps = 36000000 / conf.RatePerH
	}

	mail := NotifyMail{
		client:    c,
		store:     MailStore{conf.DB, conf.Retry},
		notifChan: nChan,
		sendChan:  make(chan M.Mail, conf.BuffLimit),
		waitTime:  conf.WaitTime,
		retry:     conf.Retry,
		closeChan: make(chan bool),
		errChan:   eChan,
		rpms:      rps,
		frpm:      conf.RetryFailedPerMin,
	}

	return &mail
}

//Set Mandrill Client
func (n *NotifyMail) SetClient(c *mandrill.Client) {
	n.client = c
}

//Get Mandrill Client
func (n *NotifyMail) GetClient() *mandrill.Client {
	return n.client
}

//Return channel for only recieving
func (n *NotifyMail) In() chan<- M.Mail {
	return n.sendChan
}

//Calls SendMail function and respondes to errors if they appear
func (n *NotifyMail) send(errChan ErrorChannel, mail M.Mail) {

	err := n.SendMail(mail)

	if err != nil {
		//If no internet connection sleep untill waitTime
		if strings.Contains(fmt.Sprintf("%s", err), "no such host") {
			time.Sleep(time.Duration(n.waitTime) * time.Second)
			errChan.In() <- errors.New("No internet connection mail")
		}

		//Try to resend mail n times
		if mail.Try >= n.retry {
			n.notifChan.In() <- fmt.Sprintf("Message sending failed after %d times %d", n.retry, mail.Id)
			mail.Status = false

			if err := n.store.Update(&mail); err != nil {
				n.errChan.In() <- err
			}

			return
		}

		//send to send mail channel again
		n.sendChan <- mail
		errChan.In() <- errors.New(fmt.Sprint("Resend: ", mail.Id))
		return
	}

	mail.Status = true
	if err := n.store.Update(&mail); err != nil {
		n.errChan.In() <- err
	}

	//If no errors send success notification
	n.notifChan.In() <- fmt.Sprintf("Mail ID: %d Succesfully sended at: %s", mail.Id)
}

//Method Uses Mandrill Api github.com/keighl/mandrill library for sending Messages
func (n *NotifyMail) SendMail(m M.Mail) error {
	message := &mandrill.Message{}

	message.AddRecipient(m.ToMail, m.ToName.String, m.ToType) // AddRecipient(toMail@gmail.com, "RecipientName", "to")
	message.FromEmail = m.FromMail                            // myemail@gmail.com
	message.FromName = m.FromName.String                      //My Name
	message.Subject = m.Subject.String                        //Subject
	message.HTML = m.Html.String                              //<h1>Long interesting HTML content</h1>

	_, apiErr, err := n.client.MessagesSend(message)

	if err != nil {
		//If Api Error ex(invalid API key, invalid Mail, ....)
		if apiErr != nil {
			return errors.New(
				fmt.Sprintf("Mandrill Api Error: Status: %s,\n Code: %d,\n Name: %s,\n Message: %s\n",
					apiErr.Status,
					apiErr.Code,
					apiErr.Name,
					apiErr.Message))
		}
		return err
	}

	return nil
}

//Starts Mails Service for listening senChannel
//And sends mails when it appears in channel
//Starts Run failed go routine for sending failed messages per n minute
func (n *NotifyMail) Start() {

	//Run separate goroutine for non blocking main program & resend failed messages per n minute.
	go func() {
		n.notifChan.In() <- fmt.Sprintf("Resend failed messages per %d minute", n.frpm)
		tick := time.NewTicker(time.Duration(n.frpm) * time.Minute)
		for {
			select {
			//Wait rate per minute tick
			case <-tick.C:
				n.ResendFailed()
			case <-n.closeChan:
				return
			}
		}
	}()

	//Starting separate goroutine for non blocking main program
	go func() {
		//Creating tick for rate per second limiting
		tick := time.NewTicker(time.Duration(n.rpms) * time.Millisecond)

		for {
			//Waiting for messages
			mail := <-n.sendChan

			//Error Channel for listening mail send errors
			errChan := NewErrorChannel(2)

			select {

			//Wait rate per second tick
			case <-tick.C:
				go n.send(errChan, mail)

			//If message sending failed send another
			case err := <-errChan.Out():
				n.errChan.In() <- err
				go n.send(errChan, mail)
			case <-n.closeChan:
				return
			}
		}
	}()
}

//Insert to database Notification and Trys to send
//If error occurs it tryes to retry n times and saves it to database as failed
func (n *NotifyMail) Notify(mail M.Mail) {
	if err := n.store.Insert(&mail); err != nil {
		n.errChan.In() <- err
		return
	}
	n.sendChan <- mail
}

//Get Failed Messages from DB and resend them
func (n *NotifyMail) ResendFailed() {

	var mails []M.Mail

	err := n.store.GetFailed(&mails, 100)

	if err != nil {
		n.errChan.In() <- err
		return
	}

	if len(mails) > 0 {
		n.notifChan.In() <- fmt.Sprintf("Resending %d failed mails", len(mails))
	}

	for _, mail := range mails {
		n.sendChan <- mail
	}

}

//Kills service at all
func (n *NotifyMail) Kill() {
	n.closeChan <- true
	close(n.sendChan)
}
