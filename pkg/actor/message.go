package actor

import "fmt"

type Message struct {
	From       *Address
	To         *Address
	Body       any
	WithReturn bool
	ReturnChan chan WrappedMessage
}

var EmptyMessage = Message{}

func NewMessage(to *Address, from *Address, body any, withReturn bool) Message {
	var c chan WrappedMessage
	if withReturn {
		c = make(chan WrappedMessage, 1)
	}
	return Message{
		To:         to,
		From:       from,
		Body:       body,
		WithReturn: withReturn,
		ReturnChan: c,
	}
}

func NewReturnMessage(body any, originalMessage Message) Message {
	return NewMessage(
		originalMessage.From,
		originalMessage.To,
		body,
		false,
	)
}

type WrappedMessage struct {
	Message *Message
	Err     error
}

func NewWrappedMessage(msg *Message, err error) WrappedMessage {
	return WrappedMessage{msg, err}
}

type AddSubscriptionMessageBody struct{}

func NewAddSubcriptionMessage(subscriberAddress *Address, notifierAddress *Address) Message {
	return Message{
		From: subscriberAddress,
		To:   notifierAddress,
		Body: AddSubscriptionMessageBody{},
	}
}

type RemoveSubscriptionMessageBody struct{}

func NewRemoveSubscriptionMessage(subscriberAddress *Address, notifierAddress *Address) Message {
	return Message{
		From: subscriberAddress,
		To:   notifierAddress,
		Body: RemoveSubscriptionMessageBody{},
	}
}

func NewSubscribersMessage(from *Address, body any) Message {
	return Message{
		From: from,
		Body: body,
	}
}

func (this *Message) String() string {
	return fmt.Sprintf("from: %s to: %s with body: %v", this.From.String(), this.To.String(), this.Body)
}
