package actor_test

import (
	"testing"
	"time"

	"github.com/pix303/cinecity/pkg/actor"
	"github.com/pix303/cinecity/pkg/actor/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func setup() (reciver *actor.Actor, sender *actor.Actor) {
	reciver = testutils.GenerateActorForTest(actor.NewAddress("lcl", "Actor-A-"))
	sender = testutils.GenerateActorForTest(actor.NewAddress("lcl", "Actor-B-"))
	return
}

func Test_ShouldActorsBeDeactivateAtCreation(t *testing.T) {
	reciver, sender := setup()

	assert.True(t, reciver.IsClosed())
	assert.True(t, sender.IsClosed())
}

func Test_ShouldActorsActivate(t *testing.T) {
	reciver, sender := setup()
	reciver.Activate()
	sender.Activate()

	assert.False(t, reciver.IsClosed())
	assert.False(t, sender.IsClosed())
	actor.Shutdown()
}

func Test_ShouldSendMessageAndChangeActorState(t *testing.T) {
	reciver, sender := setup()

	actor.RegisterActor(reciver)
	actor.RegisterActor(sender)

	var msgBody testutils.FirstMessage = "one"
	msg := actor.NewMessage(reciver.GetAddress(), sender.GetAddress(), msgBody, false)
	actor.SendMessage(msg)
	<-time.After(time.Millisecond * 10)
	state := reciver.GetState().(*string)
	assert.Contains(t, *state, "first event")
	assert.Contains(t, *state, "one")
	actor.Shutdown()
}

func Test_ShouldSendWrongMessageAndNotChangeActorState(t *testing.T) {
	reciver, sender := setup()

	actor.RegisterActor(reciver)
	actor.RegisterActor(sender)

	var msgBody testutils.WrongMessage = "wrong message"
	msg := actor.NewMessage(reciver.GetAddress(), sender.GetAddress(), msgBody, false)
	actor.SendMessage(msg)
	<-time.After(time.Millisecond * 10)
	state := reciver.GetState().(*string)
	assert.Nil(t, state)
	actor.Shutdown()
}

func Test_ShouldSendMessageAndReciveMessageBack(t *testing.T) {
	reciver, sender := setup()

	actor.RegisterActor(reciver)
	actor.RegisterActor(sender)

	var msgBody testutils.ThirdMessage = "third message"
	msg := actor.NewMessage(reciver.GetAddress(), sender.GetAddress(), msgBody, true)
	rmsg, err := actor.SendMessageWithResponse(msg)
	<-time.After(time.Millisecond * 10)
	assert.Nil(t, err)
	if rs, ok := rmsg.Body.(testutils.TestReturnMessage); ok {
		assert.Contains(t, rs, "return msg trig")
	} else {
		assert.Fail(t, "wrong message type")
	}
	actor.Shutdown()
}

// 	var secondEvent actor.SecondMessage = "two"
// 	msg := actor.NewMessage(fromPID, reciver.GetAddress(), firstMsgBody, false)
// 	reciver.Send(secondEvent, fromPID, nil)
// 	<-time.After(time.Millisecond * 10)
// 	assert.Contains(t, state.Data, "first event")

// 	reciver.Activate()

// 	assert.False(t, reciver.IsClosed())
// 	reciver.Send(secondEvent, fromPID, nil)
// 	<-time.After(time.Millisecond * 10)
// 	assert.Contains(t, state.Data, "second event")
// 	assert.Contains(t, state.Data, "two")

// 	reciver.Drop()
// 	assert.Nil(t, reciver.GetAddress())
// 	assert.Nil(t, reciver.GetMessageProcessor())

// 	actor.Shutdown()
// 	assert.Equal(t, 0, actor.NumActors())
// }
