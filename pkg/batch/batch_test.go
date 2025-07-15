package batch

import (
	"testing"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestBatcher(t *testing.T) {
	to, from := actor.GenerateAddressForTest("test")

	var msg1Body actor.FirstMessage = "hello"
	msg1 := actor.NewMessage(
		from,
		to,
		msg1Body,
		nil,
	)

	var msg2Body actor.SecondMessage = "world"
	msg2 := actor.NewMessage(
		from,
		to,
		msg2Body,
		nil,
	)

	a := actor.GenerateActorForTest("test")
	b := NewBatcher(1000, 2, a)
	a.Activate()
	b.Add(msg1)
	b.Add(msg2)

	<-time.After(100 * time.Millisecond)

	as := a.GetMessageProcessor().(*actor.TestProcessorState)
	assert.Contains(t, as.Data, "world")

	b = NewBatcher(50, 2, a)
	a.Activate()
	b.Add(msg1)
	<-time.After(100 * time.Millisecond)
	b.Add(msg2)
	as = a.GetMessageProcessor().(*actor.TestProcessorState)
	assert.Contains(t, as.Data, "hello")

	<-time.After(100 * time.Millisecond)
	assert.Contains(t, as.Data, "world")

}
