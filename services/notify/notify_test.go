package notify

import (
	"github.com/keighl/mandrill"
	"testing"
	"time"
)

//Creates new Config instantce for testing
func newConf() Conf {
	ntfChan = NewNotificationChannel(10)
	errChan = NewErrorChannel(10)
	wt := 1
	tm := 1

	mailConf = NotifyMailConfig{"Apikey", db, wt, tm, 100, 3, 250, 1}
	return Conf{
		MailConf:  mailConf,
		BuffLimit: 100,
	}
}

func TestNotificationServiceMailSend(t *testing.T) {

	nfs := NewNotificationService(newConf())

	server, cl := newFakeServer(200, `[{"email":"ahidoyatov@gmail.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()

	nfs.MailService.SetClient(&mandrill.Client{"APIKEY", server.URL, cl})

	nfs.Start()

	mail := newMail(906)
	nfs.Send(mail)

	select {

	case <-nfs.Notifications():
	case <-nfs.Errors():
		t.Error("Unexpected Error")
	case <-time.After(1 * time.Second):
		t.Error("Mail Send Error")

	}
}
