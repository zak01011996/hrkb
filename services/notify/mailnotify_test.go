package notify

import (
	"github.com/keighl/mandrill"
	"testing"
	"time"
)

var (
	mailntf  *NotifyMail
	ntfChan  NotificationChannel
	mailConf NotifyMailConfig
	errChan  ErrorChannel
)

func init() {
	ntfChan = NewNotificationChannel(10)
	errChan = NewErrorChannel(10)
	wt := 1
	tm := 1

	mailConf = NotifyMailConfig{"Apikey", db, wt, tm, 10, 3, 250, 2}
	mailntf = NewNotifyMail(mailConf, ntfChan, errChan)
}

func TestMessageSendSuccess(t *testing.T) {
	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()

	mail := newMail(900)

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	errChan := NewErrorChannel(4)

	go mailntf.send(errChan, mail)

	select {
	case <-mailntf.notifChan.Out():
	case <-mailntf.sendChan:
		t.Error("Mail Send Error")
	case <-time.After(1 * time.Second):
		t.Error("Mail Send Error")
	}
}

func TestMessageSendFail(t *testing.T) {
	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	server.Close()

	mail := newMail(901)
	mail.Try = mailConf.Retry

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	errChan := NewErrorChannel(4)
	go mailntf.send(errChan, mail)

	select {
	case <-ntfChan.Out():
	case <-time.After(1 * time.Second):
		t.Error("Expected error")
	}
}

func TestMessageSendRetry(t *testing.T) {
	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	server.Close()

	mail := newMail(901)

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	errChan := NewErrorChannel(4)

	go mailntf.send(errChan, mail)

	select {
	case <-mailntf.sendChan:
	case <-time.After(1 * time.Second):
		t.Error("Retry error")
	}
}

func TestMessageSendTimeout(t *testing.T) {
	server, cl := newSleepyServer(500, 200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()

	mail := newMail(901)

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	errChan := NewErrorChannel(4)
	go mailntf.send(errChan, mail)

	select {
	case <-mailntf.sendChan:
	case <-time.After(2 * time.Second):
		t.Error("Timeout error")
	}
}

func TestMailNotifyStart(t *testing.T) {
	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()

	mail := newMail(905)

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	mailntf.Start()

	mailntf.Notify(mail)

	select {
	//Wait to Notify Mail Start
	case <-ntfChan.Out():
		//Wait to Resend Failed start
		<-ntfChan.Out()
	case <-time.After(3 * time.Second):
		t.Error("Mail notify Start Error")
	}
}

func TestMailNotifyResendFailed(t *testing.T) {
	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()

	mailntf.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	//Insert to database failed tryed 1 time message for getting at least 1 failed message
	mail := newMail(901)

	mail.Status = false
	mail.Try = 1
	mail.Created = time.Now()
	mail.Active = true

	if _, err := mailntf.store.db.Insert(&mail); err != nil {
		t.Error("Db insert error")
	}

	mailntf.Start()
	mailntf.ResendFailed()

	select {
	case <-mailntf.sendChan:
	case <-time.After(1 * time.Second):
		t.Error("Mail notify ResendFailed Error")
	}
}
