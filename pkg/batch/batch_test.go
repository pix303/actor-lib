package batch_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pix303/cinecity/pkg/actor"
	"github.com/pix303/cinecity/pkg/batch"
	"github.com/stretchr/testify/assert"
)

func setup() (msg1 actor.Message, msg2 actor.Message, actor1 *actor.Actor, handler func(msg actor.Message)) {

	to, from := actor.GenerateAddressForTest("test")

	var msg1Body actor.FirstMessage = "hello"
	msg1 = actor.NewMessage(
		from,
		to,
		msg1Body,
		nil,
	)

	var msg2Body actor.SecondMessage = "world"
	msg2 = actor.NewMessage(
		from,
		to,
		msg2Body,
		nil,
	)

	actor1 = actor.GenerateActorForTest("test")
	actor1.Activate()
	handler = func(msg actor.Message) {
		s := actor1.GetMessageProcessor().(*actor.TestProcessorState)
		if v, ok := msg.Body.(actor.FirstMessage); ok {
			s.Data = fmt.Sprintf("write by batcher: %s", v)
		}
	}
	return
}

func TestBatcher_StateChangesAfterTimeoutLimits(t *testing.T) {
	msg1, msg2, a, handler := setup()

	b := batch.NewBatcher(1000, 100, handler)
	b.Add(msg1)
	b.Add(msg2)

	as := a.GetMessageProcessor().(*actor.TestProcessorState)
	<-time.After(1001 * time.Millisecond)

	assert.Contains(t, as.Data, "hello")
}

func TestBatcher_StateDontChangesBeforeTimeoutLimits(t *testing.T) {
	msg1, msg2, a, handler := setup()

	b := batch.NewBatcher(1000, 100, handler)
	b.Add(msg1)
	b.Add(msg2)

	<-time.After(100 * time.Millisecond)

	as := a.GetMessageProcessor().(*actor.TestProcessorState)
	assert.NotContains(t, as.Data, "hello")
}

func TestBatcher_StateChangesByMaxItemsLimits(t *testing.T) {
	msg1, msg2, a, handler := setup()

	b := batch.NewBatcher(1000, 2, handler)
	b.Add(msg1)
	b.Add(msg2)

	<-time.After(100 * time.Millisecond)

	as := a.GetMessageProcessor().(*actor.TestProcessorState)
	assert.Contains(t, as.Data, "hello")
}

func TestBatcher_StateDontChangesBeforMaxItemsLimits(t *testing.T) {
	msg1, msg2, a, handler := setup()

	b := batch.NewBatcher(1000, 3, handler)
	b.Add(msg1)
	b.Add(msg2)

	<-time.After(999 * time.Millisecond)

	as := a.GetMessageProcessor().(*actor.TestProcessorState)
	assert.NotContains(t, as.Data, "hello")
}
