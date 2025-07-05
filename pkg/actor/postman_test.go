package actor_test

import (
	"testing"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	actor.Shutdown()
	p := actor.GetPostman()
	assert.NotNil(t, p)

	a := actor.GenerateActorForTest("A")
	b := actor.GenerateActorForTest("B")
	a.Activate()
	b.Activate()
	assert.Contains(t, a.GetAddress().String(), "A-sender")
	assert.Contains(t, b.GetAddress().String(), "B-sender")

	actor.RegisterActor(a)
	actor.RegisterActor(b)
	assert.Equal(t, 2, actor.NumActors())

	var re actor.ThirdMessage = "three"
	msg := actor.Message{
		From: a.GetAddress(),
		To:   b.GetAddress(),
		Body: re,
	}

	a.Send(msg)
	<-time.After(time.Millisecond * 300)

	astate := a.GetMessageProcessor().(*actor.TestProcessorState)
	bstate := b.GetMessageProcessor().(*actor.TestProcessorState)

	assert.Contains(t, bstate.Data, "three")
	assert.Contains(t, astate.Data, "return msg")

	var returnMsgBody actor.WithSyncResponse = "wait response"
	rm := actor.Message{
		From: a.GetAddress(),
		To:   a.GetAddress(),
		Body: returnMsgBody,
	}
	returnMsg, err := actor.DispatchMessageWithReturn(rm)
	assert.Nil(t, err)
	assert.Contains(t, returnMsg.Body, "message recived")

	assert.Equal(t, 2, actor.NumActors())
	actor.Shutdown()
	assert.Equal(t, 0, actor.NumActors())
}
