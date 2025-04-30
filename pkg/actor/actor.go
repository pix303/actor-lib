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
	address    *Address
	inbox      chan Message
	withReturn bool
	isClosed   bool
	processor  Processor
	postman    *Postman
}

func NewActor(address *Address, processor Processor) (Actor, error) {
	inbox := make(chan Message, 100)
	if address == nil {
		return Actor{}, AddressNilErr
	}
	return Actor{
		address:   address,
		processor: processor,
		inbox:     inbox,
		isClosed:  true,
	}, nil
}

func (this *Actor) Activate() {
	if this.IsClosed() {
		slog.Info("Actor activated", slog.String("address", this.address.String()))
		this.isClosed = false
		p := this.GetProcessor()
		if p != nil {
			go p.Process(this.inbox)
		}
	}
}

func (this *Actor) GetAddress() *Address {
	return this.address
}

func (this *Actor) IsClosed() bool {
	return this.isClosed
}

func (this *Actor) GetProcessor() Processor {
	return this.processor
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
	this.inbox <- msg
	return nil
}

func (this *Actor) Send(msg Message) {
	DispatchMessage(msg)
}

func (this *Actor) Drop() {
	this.Deactivate()
	this.address = nil
	this.processor = nil
	this.inbox = nil
	this.postman = nil
}

func (this *Actor) String() string {
	return fmt.Sprintf("address: %s - isClosed: %t", this.address.String(), this.IsClosed())
}
