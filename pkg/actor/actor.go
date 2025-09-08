package actor

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrAddressNil  = errors.New("actor cant have a nil address")
	ErrInboxClosed = errors.New("actor has inbox closed")
)

type Actor struct {
	address          *Address
	MessageBox       chan Message
	isClosed         bool
	messageProcessor MessageProcessor
	postman          *Postman
}

func NewActor(address *Address, processor MessageProcessor) (Actor, error) {
	if address == nil {
		return Actor{}, ErrAddressNil
	}
	return Actor{
		address:          address,
		messageProcessor: processor,
		MessageBox:       make(chan Message, 100),
		isClosed:         true,
	}, nil
}

func (actor *Actor) Activate() {
	if actor.IsClosed() {
		slog.Info("Actor activated", slog.String("address", actor.address.String()))
		actor.isClosed = false
		p := actor.GetMessageProcessor()
		if p != nil {
			go p.Process(actor.MessageBox)
		}
	}
}

func (actor *Actor) GetAddress() *Address {
	return actor.address
}

func (actor *Actor) IsClosed() bool {
	return actor.isClosed
}

func (actor *Actor) GetMessageProcessor() MessageProcessor {
	return actor.messageProcessor
}

func (actor *Actor) Deactivate() {
	if !actor.IsClosed() {
		actor.isClosed = true
		slog.Info("Actor deactivated", slog.String("address", actor.address.String()))
	}
}

func (actor *Actor) Inbox(msg Message) error {
	if actor.isClosed {
		return ErrInboxClosed
	}
	actor.MessageBox <- msg
	return nil
}

func (actor *Actor) Send(msg Message) error {
	return SendMessage(msg)
}

func (actor *Actor) Drop() {
	mp := actor.GetMessageProcessor()
	if mp != nil {
		mp.Shutdown()
	}
	actor.Deactivate()
	actor.address = nil
	actor.messageProcessor = nil
	actor.MessageBox = nil
	actor.postman = nil
}

func (actor *Actor) String() string {
	return fmt.Sprintf("address: %s - isClosed: %t", actor.address.String(), actor.IsClosed())
}
