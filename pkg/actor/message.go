package actor

import "fmt"

type Message struct {
	From       *Address
	To         *Address
	Body       any
	WithReturn chan<- Message
}

func EmptyMessage() Message {
	return Message{}
}

func NewMessage(to *Address, from *Address, body any, withReturn chan<- Message) Message {
	return Message{
		To:         to,
		From:       from,
		Body:       body,
		WithReturn: withReturn,
	}
}

func (this *Message) String() string {
	return fmt.Sprintf("from: %s to: %s with body: %v", this.From.String(), this.To.String(), this.Body)
}
