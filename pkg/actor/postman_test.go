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
	msg := actor.NewMessage(
		b.GetAddress(),
		a.GetAddress(),
		re,
		nil,
	)

	a.Send(msg)
	<-time.After(time.Millisecond * 100)

	astate := a.GetMessageProcessor().(*actor.TestProcessorState)
	bstate := b.GetMessageProcessor().(*actor.TestProcessorState)

	assert.Contains(t, bstate.Data, "three")
	assert.Contains(t, astate.Data, "return msg")

	var returnMsgBody actor.WithSyncResponse = "wait response"
	withReturnMessage := actor.NewMessage(
		b.GetAddress(),
		a.GetAddress(),
		returnMsgBody,
		a.MessageBox,
	)

	a.Send(withReturnMessage)
	<-time.After(time.Millisecond * 200)

	astate = a.GetMessageProcessor().(*actor.TestProcessorState)
	bstate = b.GetMessageProcessor().(*actor.TestProcessorState)
	assert.Contains(t, bstate.Data, "wait response")
	assert.Contains(t, astate.Data, "message recived")

	assert.Equal(t, 2, actor.NumActors())
	actor.Shutdown()
	assert.Equal(t, 0, actor.NumActors())
}
