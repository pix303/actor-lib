package actor

import (
	"fmt"
	"log/slog"
)

type TestProcessorState struct {
	Data string
}

type FirstMessage string
type SecondMessage string
type ThirdMessage string
type ReturnMessage string
type WithSyncResponse string
type Response string

func (this *TestProcessorState) Process(inbox chan Message) {
	for {
		msg := <-inbox
		switch msg.Body.(type) {
		case FirstMessage:
			this.Data = fmt.Sprintf("processed with first event: %s", msg.Body)
		case SecondMessage:
			this.Data = fmt.Sprintf("processed with second event: %s", msg.Body)
		case ThirdMessage:
			this.Data = fmt.Sprintf("processed with third event: %s", msg.Body)
			var r ReturnMessage = "return msg"
			rmsg := Message{
				From: msg.To,
				To:   msg.From,
				Body: r,
			}
			DispatchMessage(rmsg)
		case ReturnMessage:
			this.Data = fmt.Sprintf("processed with return event: %s", msg.Body)
		}
	}
}

func (this *TestProcessorState) ProcessSync(msg Message) Message {

	switch msg.Body.(type) {
	case WithSyncResponse:
		this.Data = fmt.Sprintf("processed with sync message: %s", msg.Body)
		var rm = Message{
			To:   msg.From,
			From: *NewAddress("local", "me"),
			Body: "message recived",
		}
		return rm

	default:
		var em = Message{
			To:   msg.From,
			From: *NewAddress("local", "me"),
			Body: "empty message",
		}
		return em
	}
}

func (this *TestProcessorState) Shutdown() {
	this.Data = ""
	slog.Info("all clean after shutdown")
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
