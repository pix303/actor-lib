package batch

import (
	"log/slog"
	"sync"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
)

type Batcher struct {
	messages       []actor.Message
	maxNumMessages uint
	timeout        time.Duration
	timer          *time.Timer
	mutex          sync.Mutex
	actorReciver   *actor.Actor
}

func NewBatcher(timeoutMs uint, maxMessages uint, actorReciver *actor.Actor) *Batcher {
	b := Batcher{
		timeout:        time.Duration(timeoutMs) * time.Millisecond,
		messages:       make([]actor.Message, 0),
		mutex:          sync.Mutex{},
		maxNumMessages: maxMessages,
		actorReciver:   actorReciver,
	}
	slog.Info("Batcher created", "timeout", timeoutMs, "maxMessages", maxMessages)
	return &b
}

func (this *Batcher) Add(msg actor.Message) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.messages) == 0 {
		this.timer = time.AfterFunc(this.timeout, this.process)
	}

	this.messages = append(this.messages, msg)
	if len(this.messages) == int(this.maxNumMessages) {
		this.process()
	}
}

func (this *Batcher) process() {
	for _, x := range this.messages {
		this.actorReciver.Inbox(x)
	}

	this.timer.Stop()
	this.messages = make([]actor.Message, 0)
	slog.Info("batch process end")
}
