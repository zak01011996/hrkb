package notify

import (
	"errors"
	M "hrkb/models"
)

type Conf struct {
	MailConf  NotifyMailConfig
	BuffLimit int
	TgGroupId int
}

type NotificationService struct {
	ntfChan     NotificationChannel
	errChan     ErrorChannel
	sendChan    SendChannel
	closeChan   chan bool
	MailService *NotifyMail
	TgService   *NotifyTelegram
}

//Return new instance of Notification Service
func NewNotificationService(c Conf) *NotificationService {

	ntfChan := NewNotificationChannel(c.BuffLimit + 10)
	errChan := NewErrorChannel(c.BuffLimit + 50)
	sendChan := NewSendChannel(c.BuffLimit)

	nfs := NotificationService{
		ntfChan:     ntfChan,
		errChan:     errChan,
		sendChan:    sendChan,
		closeChan:   make(chan bool),
		MailService: NewNotifyMail(c.MailConf, ntfChan, errChan),
		TgService:   NewNotifyTelegram(c.BuffLimit, c.TgGroupId, ntfChan, errChan),
	}

	return &nfs
}

//Starts Notification service then child services
func (n *NotificationService) Start() {
	n.MailService.Start()
	n.TgService.Start()

	n.ntfChan.In() <- "Notification Service started"

	//Starting separate go routine for not blocking main programm
	//Go routine listens to send Channel. When interface comes to channel type casts it and if type true sends
	//it proper service else returns error
	go func() {
		for {
			select {
			case obj := <-n.sendChan.Out():
				//Type casting
				switch o := obj.(type) {
				case M.Mail:
					n.MailService.Notify(o)
				case string:
					n.TgService.Notify(o)
				default:
					n.errChan.In() <- errors.New("Type Cast error")
				}

			case <-n.closeChan:
				return
			}
		}
	}()
}

//Send Notification interface to Send Channel
func (n *NotificationService) Send(i interface{}) {
	n.sendChan.In() <- i
}

//Resend failed messages
func (n *NotificationService) ResendFailed() {
	n.MailService.ResendFailed()
}

//Returns Channel For listening notifications
func (n *NotificationService) Notifications() <-chan string {
	return n.ntfChan.Out()
}

//Returns Channel For listening errors
func (n *NotificationService) Errors() <-chan error {
	return n.errChan.Out()
}

//Sends Stop signal to Services, then closes Send Channel
//Kills service at all
func (n *NotificationService) Kill() {
	n.closeChan <- true
	close(n.sendChan)
}
