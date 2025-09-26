package actor_test

import (
	"testing"
	"time"

	"github.com/pix303/cinecity/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	actor.Shutdown()
	p := actor.GetPostman()
	assert.NotNil(t, p)

	a := GenerateActorForTest("A")
	b := GenerateActorForTest("B")
	a.Activate()
	b.Activate()
	assert.Contains(t, a.GetAddress().String(), "A-sender")
	assert.Contains(t, b.GetAddress().String(), "B-sender")

	actor.RegisterActor(a)
	actor.RegisterActor(b)
	assert.Equal(t, 2, actor.NumActors())

	var re ThirdMessage = "three"
	a.Send(re, a.GetAddress(), nil)
	<-time.After(time.Millisecond * 100)

	astate := a.GetMessageProcessor().(*TestProcessorState)
	bstate := b.GetMessageProcessor().(*TestProcessorState)

	assert.Contains(t, bstate.Data, "three")
	assert.Contains(t, astate.Data, "return msg")

	var returnMsgBody WithSyncResponse = "wait response"
	a.InboxAndWaitResponse(returnMsgBody, b.GetAddress())
	<-time.After(time.Millisecond * 200)

	astate = a.GetMessageProcessor().(*TestProcessorState)
	bstate = b.GetMessageProcessor().(*TestProcessorState)
	assert.Contains(t, bstate.Data, "wait response")
	assert.Contains(t, astate.Data, "message recived")

	assert.Equal(t, 2, actor.NumActors())
	actor.Shutdown()
	assert.Equal(t, 0, actor.NumActors())
}
