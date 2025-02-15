package actor_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

type TestProcessor struct {
	data string
}

func (this *TestProcessor) Process(inbox chan actor.Message) {
	for msg := range inbox {
		this.data = fmt.Sprintf("processed by %s", msg.ToString())
	}
}

func GeneratePIDs(prefix string) (toPID *actor.PID, fromPID *actor.PID) {
	if prefix == "" {
		prefix = "test"
	}
	toPID = actor.NewPID("local", prefix+"-sender")
	fromPID = actor.NewPID("local", prefix+"-reciver")
	return
}

func GenerateActor(prefix string) *actor.Actor {
	processor := TestProcessor{}
	toPID, _ := GeneratePIDs(prefix)
	a := actor.NewActor(toPID, &processor)
	return &a
}

func TestActor(t *testing.T) {
	toPID, fromPID := GeneratePIDs("test")
	assert.Equal(t, toPID.String(), "local.test-sender")
	a := GenerateActor("test")
	processor := a.GetProcessor().(*TestProcessor)

	assert.True(t, a.IsClosed())
	a.Activate()
	assert.False(t, a.IsClosed())

	a.SendToMe(actor.Message{To: *toPID, From: *fromPID, Body: "hello"})
	<-time.After(time.Millisecond * 100)
	assert.Contains(t, a.GetAddress().String(), "test-sender")
	assert.Contains(t, processor.data, "test-sender")
	assert.Contains(t, processor.data, "test-reciver")
	assert.Contains(t, processor.data, "hello")

	a.Deactivate()

	assert.True(t, a.IsClosed())
	a.SendToMe(actor.Message{To: *toPID, From: *fromPID, Body: "deactivate test processor state must not change"})
	<-time.After(time.Millisecond * 100)
	assert.Contains(t, processor.data, "hello")
	assert.NotContains(t, processor.data, "deactivate")
	assert.NotContains(t, processor.data, "must not change")

	a.Activate()

	assert.False(t, a.IsClosed())
	a.SendToMe(actor.Message{To: *toPID, From: *fromPID, Body: "reactivate test processor state must change"})
	<-time.After(time.Millisecond * 100)
	assert.Contains(t, processor.data, "reactivate test")
	assert.Contains(t, processor.data, "must change")

	a.Drop()
	assert.Nil(t, a.GetAddress())
	assert.Nil(t, a.GetProcessor())
}
