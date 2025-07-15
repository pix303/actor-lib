package actor

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	AddressNilErr  = errors.New("actor cant have a nil address")
	InboxClosedErr = errors.New("actor has inbox closed")
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
		return Actor{}, AddressNilErr
	}
	return Actor{
		address:          address,
		messageProcessor: processor,
		MessageBox:       make(chan Message, 100),
		isClosed:         true,
	}, nil
}

func (this *Actor) Activate() {
	if this.IsClosed() {
		slog.Info("Actor activated", slog.String("address", this.address.String()))
		this.isClosed = false
		p := this.GetMessageProcessor()
		if p != nil {
			go p.Process(this.MessageBox)
		}
	}
}

func (this *Actor) GetAddress() *Address {
	return this.address
}

func (this *Actor) IsClosed() bool {
	return this.isClosed
}

func (this *Actor) GetMessageProcessor() MessageProcessor {
	return this.messageProcessor
}

func (this *Actor) Deactivate() {
	if this.IsClosed() == false {
		this.isClosed = true
		slog.Info("Actor deactivated", slog.String("address", this.address.String()))
	}
}

func (this *Actor) Inbox(msg Message) error {
	if this.isClosed {
		return InboxClosedErr
	}
	this.MessageBox <- msg
	return nil
}

func (this *Actor) Send(msg Message) {
	SendMessage(msg)
}

func (this *Actor) Drop() {
	mp := this.GetMessageProcessor()
	if mp != nil {
		mp.Shutdown()
	}
	this.Deactivate()
	this.address = nil
	this.messageProcessor = nil
	this.MessageBox = nil
	this.postman = nil
}

func (this *Actor) String() string {
	return fmt.Sprintf("address: %s - isClosed: %t", this.address.String(), this.IsClosed())
}
