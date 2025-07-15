package actor_test

import (
	"testing"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestActor(t *testing.T) {
	toPID, fromPID := actor.GenerateAddressForTest("test")

	a := actor.GenerateActorForTest("test")
	state := a.GetMessageProcessor().(*actor.TestProcessorState)

	assert.True(t, a.IsClosed())

	actor.RegisterActor(a)

	var firstEvent actor.FirstMessage = "one"
	a.Send(actor.Message{To: toPID, From: fromPID, Body: firstEvent})
	<-time.After(time.Millisecond * 10)
	assert.Contains(t, a.GetAddress().String(), "test-sender")
	assert.Contains(t, state.Data, "first event")
	assert.Contains(t, state.Data, "one")

	a.Deactivate()

	assert.True(t, a.IsClosed())
	var secondEvent actor.SecondMessage = "two"
	a.Send(actor.Message{To: toPID, From: fromPID, Body: secondEvent})
	<-time.After(time.Millisecond * 10)
	assert.Contains(t, state.Data, "first event")

	a.Activate()

	assert.False(t, a.IsClosed())
	a.Send(actor.Message{To: toPID, From: fromPID, Body: secondEvent})
	<-time.After(time.Millisecond * 10)
	assert.Contains(t, state.Data, "second event")
	assert.Contains(t, state.Data, "two")

	a.Drop()
	assert.Nil(t, a.GetAddress())
	assert.Nil(t, a.GetMessageProcessor())

	actor.Shutdown()
	assert.Equal(t, 0, actor.NumActors())
}
