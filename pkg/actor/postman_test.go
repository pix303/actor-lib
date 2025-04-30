package actor_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	slog.Info("---- start test dispa")
	actor.DropAllActors()
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

	var re actor.ThirdEvent = "three"
	msg := actor.Message{
		From: *a.GetAddress(),
		To:   *b.GetAddress(),
		Body: re,
	}

	a.Send(msg)
	<-time.After(time.Millisecond * 300)

	astate := a.GetProcessor().(*actor.TestProcessorState)
	bstate := b.GetProcessor().(*actor.TestProcessorState)

	assert.Contains(t, bstate.Data, "three")
	assert.Contains(t, astate.Data, "return msg")

	assert.Equal(t, actor.NumActors(), 2)
	actor.DropAllActors()
	assert.Equal(t, actor.NumActors(), 0)
	slog.Info("---- end test dispa")
}
