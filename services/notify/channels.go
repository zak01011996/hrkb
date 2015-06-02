package notify

import ()

type Channel interface {
	In() chan<- interface{}  // Should return channel for only sending
	Out() <-chan interface{} // Should return channel for only getting
	Close()                  //Close channel
}

type ErrChannel interface {
	In() chan<- error  // Should return channel for only sending
	Out() <-chan error // Should return channel for only getting
	Close()            //Close channel
}

//Channel for sending notifications
type SendChannel chan interface{}

func NewSendChannel(size int) SendChannel {
	return make(chan interface{}, size)
}

func (ch SendChannel) In() chan<- interface{} {
	return ch
}

func (ch SendChannel) Out() <-chan interface{} {
	return ch
}

func (ch SendChannel) Close() {
	close(ch)
}

//Channel for listening outcoming notifications
type NotificationChannel chan string

func NewNotificationChannel(size int) NotificationChannel {
	return make(chan string, size)
}

func (ch NotificationChannel) In() chan<- string {
	return ch
}

func (ch NotificationChannel) Out() <-chan string {
	return ch
}

func (ch NotificationChannel) Close() {
	close(ch)
}

//Channel for Sending errors
type ErrorChannel chan error

func NewErrorChannel(size int) ErrorChannel {
	return make(chan error, size)
}

func (ch ErrorChannel) In() chan<- error {
	return ch
}

func (ch ErrorChannel) Out() <-chan error {
	return ch
}

func (ch ErrorChannel) Close() {
	close(ch)
}
