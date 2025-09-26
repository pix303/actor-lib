package testutils

import (
	"fmt"
	"github.com/pix303/cinecity/pkg/actor"
	"log/slog"
)

type TestProcessorState struct {
	Data *string
}

type WrongMessage string
type FirstMessage string
type SecondMessage string
type ThirdMessage string
type TestReturnMessage string
type WithSyncResponse string
type Response string

func (state *TestProcessorState) GetState() any {
	return state.Data
}

func (state *TestProcessorState) Process(msg actor.Message) {
	slog.Info("processing msg", slog.String("msg", msg.String()))
	switch msg.Body.(type) {
	case FirstMessage:
		r := fmt.Sprintf("processed by first event: %s", msg.Body)
		state.Data = &r
	case SecondMessage:
		r := fmt.Sprintf("%s ++ processed by second event: %s", *state.Data, msg.Body)
		state.Data = &r
	case ThirdMessage:
		r := fmt.Sprintf("processed by third event: %s", msg.Body)
		state.Data = &r
		var returnMsg TestReturnMessage = "return msg triggerd by third message"
		rmsg := actor.Message{
			From: msg.To,
			To:   msg.From,
			Body: returnMsg,
		}
		msg.ReturnChan <- actor.NewWrappedMessage(&rmsg, nil)
	case TestReturnMessage:
		r := fmt.Sprintf("processed with return event: %s", msg.Body)
		state.Data = &r
	}
}

func (state *TestProcessorState) Shutdown() {
	state.Data = nil
	slog.Info("all clean after shutdown")
}

func GenerateActorForTest(address *actor.Address) *actor.Actor {
	processor := TestProcessorState{}
	a, _ := actor.NewActor(address, &processor)
	return &a
}
