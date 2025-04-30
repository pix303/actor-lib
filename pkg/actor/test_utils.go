package actor

import (
	"fmt"
	"log/slog"
)

type TestProcessorState struct {
	Data string
}

type FirstEvent string
type SecondEvent string
type ThirdEvent string
type ReturnEvent string

func (this *TestProcessorState) Process(inbox chan Message) {
	for {
		msg := <-inbox
		switch msg.Body.(type) {
		case FirstEvent:
			this.Data = fmt.Sprintf("processed with first event: %s", msg.Body)
		case SecondEvent:
			this.Data = fmt.Sprintf("processed with second event: %s", msg.Body)
		case ThirdEvent:
			this.Data = fmt.Sprintf("processed with third event: %s", msg.Body)
			var r ReturnEvent = "return msg"
			rmsg := Message{
				From: msg.To,
				To:   msg.From,
				Body: r,
			}
			DispatchMessage(rmsg)
		case ReturnEvent:
			this.Data = fmt.Sprintf("processed with return event: %s", msg.Body)
		}
	}
}

func GenerateAddressForTest(prefix string) (toPID *Address, fromPID *Address) {
	if prefix == "" {
		prefix = "test"
	}
	toPID = NewAddress("local", prefix+"-sender")
	fromPID = NewAddress("local", prefix+"-reciver")
	return
}

func GenerateActorForTest(prefix string) *Actor {
	processor := TestProcessorState{}
	toPID, _ := GenerateAddressForTest(prefix)
	a, _ := NewActor(toPID, &processor)
	slog.Info("generate", slog.String("a", a.GetAddress().String()))
	return &a
}
