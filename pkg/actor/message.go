package actor

import "fmt"

type Message struct {
	From Address
	To   Address
	Body any
}

func (this *Message) String() string {
	return fmt.Sprintf("from: %s to: %s with body: %v", this.From.String(), this.To.String(), this.Body)
}
