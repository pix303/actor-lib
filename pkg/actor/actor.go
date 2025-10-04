package actor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrAddressNil            = errors.New("actor cant have a nil address")
	ErrAddressAlreadyUsed    = errors.New("address is already used")
	ErrInboxClosed           = errors.New("actor has inbox closed")
	ErrSendWithReturnTimeout = errors.New("message with retrun has not be processed in time")
)

type Actor struct {
	address    *Address
	MessageBox chan Message
	isClosed   bool
	state      StateProcessor
}

func NewActor(address *Address, processor StateProcessor) (Actor, error) {
	if address == nil {
		return Actor{}, ErrAddressNil
	}

	p := GetPostman()
	if temp := p.actors[address.String()]; temp != nil {
		slog.Error(ErrAddressAlreadyUsed.Error(), slog.String("actor-address", address.String()))
		return Actor{}, ErrAddressAlreadyUsed
	}

	return Actor{
		address:    address,
		state:      processor,
		MessageBox: make(chan Message, 100),
		isClosed:   true,
	}, nil
}

func (a *Actor) Activate() {
	if a.IsClosed() {
		slog.Info("Actor activated", slog.String("address", a.address.String()))
		a.isClosed = false
		p := a.state
		if p != nil {
			go a.processMessage(a.MessageBox)
		}
	}
}

func (a *Actor) processMessage(inboxChan <-chan Message) {
	for {
		msg := <-inboxChan
		a.state.Process(msg)
	}
}

func (a *Actor) GetAddress() *Address {
	return a.address
}

func (a *Actor) IsClosed() bool {
	return a.isClosed
}

func (a *Actor) Deactivate() {
	if !a.IsClosed() {
		a.isClosed = true
		slog.Info("Actor deactivated", slog.String("address", a.address.String()))
	}
}

func (a *Actor) Inbox(msg Message) error {
	if a.isClosed {
		return ErrInboxClosed
	}
	a.MessageBox <- msg
	return nil
}

func (a *Actor) InboxAndWaitResponse(msg Message) (Message, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(msg.ReturnTimeout)*time.Second)
	defer cancelFunc()

	returnChan := msg.ReturnChan

	err := a.Inbox(msg)
	if err != nil {
		return Message{}, err
	}

	select {
	case returnMsg := <-returnChan:
		return *returnMsg.Message, returnMsg.Err
	case <-ctx.Done():
		return Message{}, ErrSendWithReturnTimeout
	}
}

func (a *Actor) Drop() {
	mp := a.state
	if mp != nil {
		mp.Shutdown()
	}
	a.Deactivate()
	UnRegisterActor(a.address)
	a.address = nil
	a.state = nil
	a.MessageBox = nil
}

func (a *Actor) GetState() any {
	mp := a.state
	if mp != nil {
		return mp.GetState()
	}
	return nil
}

func (a *Actor) String() string {
	return fmt.Sprintf("address: %s - isClosed: %t", a.address.String(), a.IsClosed())
}
