package actor

import "fmt"

type Message struct {
	From PID
	To   PID
	Body any
}

func (this *Message) ToString() string {
	return fmt.Sprintf("from: %s to: %s with body: %v", this.From.String(), this.To.String(), this.Body)
}
