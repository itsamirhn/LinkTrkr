package pkg

import "github.com/sirupsen/logrus"

const (
	SubscriberBufferSize = 10
	MessageBufferSize    = 100
)

type notifier[T any] struct {
	subscribers chan chan T
	messages    chan T
}

type Notifier[T any] interface {
	Subscribe() <-chan T
	Notify(msg T)
	Run()
}

func NewNotifier[T any]() Notifier[T] {
	return &notifier[T]{
		subscribers: make(chan chan T, MessageBufferSize),
		messages:    make(chan T, MessageBufferSize),
	}
}

func (n *notifier[T]) Subscribe() <-chan T {
	subscriber := make(chan T, SubscriberBufferSize)
	n.subscribers <- subscriber
	return subscriber
}

func (n *notifier[T]) Notify(msg T) {
	n.messages <- msg
}

func (n *notifier[T]) Run() {
	subs := make([]chan T, 0)

	for {
		select {
		case subscriber := <-n.subscribers:
			subs = append(subs, subscriber)

		case message := <-n.messages:
			for _, sub := range subs {
				select {
				case sub <- message:
				default:
					logrus.WithField("message", message).Warn("Subscriber channel full, skipping message.")
				}
			}
		}
	}
}
