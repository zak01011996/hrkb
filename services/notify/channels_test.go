package notify

import (
	"errors"
	"fmt"
	"testing"
)

func TestSendChannel(t *testing.T) {
	ch := NewSendChannel(2)

	ch.In() <- "test"
	ch.In() <- "test"

	if out := <-ch.Out(); out != "test" {
		t.Error("SendChannel recieve error")
	}
}

func TestNotificationChannel(t *testing.T) {
	ch := NewNotificationChannel(3)

	ch.In() <- "test"
	ch.In() <- "test"

	if out := <-ch.Out(); out != "test" {
		t.Error("SendChannel recieve error")
	}
}

func TestErrorChannel(t *testing.T) {
	ch := NewErrorChannel(3)

	ch.In() <- errors.New("test")
	ch.In() <- errors.New("test")

	if out := <-ch.Out(); fmt.Sprintf("%s", out) != "test" {
		t.Error("SendChannel recieve error")
	}
}
