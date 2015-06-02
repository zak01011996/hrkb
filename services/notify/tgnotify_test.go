package notify

import (
	"testing"
	"time"
)

func TestNotifyTelegramNotify(t *testing.T) {

	ntfChan := NewNotificationChannel(10)
	errChan := NewErrorChannel(10)
	ntg := NewNotifyTelegram(100, 17006774, ntfChan, errChan)

	ntg.Start()

	ntg.Notify("Test")

	select {
	case <-ntfChan.Out():
	case <-errChan.Out():
		t.Error("Notification send error")
	case <-time.After(2 * time.Second):
		t.Error("Notification send error")
	}
}
