package actor

import "fmt"

type Message struct {
	From       *Address
	To         *Address
	Body       any
	WithReturn chan<- Message
}

var EmptyMessage = Message{}

func NewMessage(to *Address, from *Address, body any, withReturn chan<- Message) Message {
	return Message{
		To:         to,
		From:       from,
		Body:       body,
		WithReturn: withReturn,
	}
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
