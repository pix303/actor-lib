package actor

import "log/slog"

type Dispatchable interface {
	GetAddress() *PID
	IsClosed() bool
	Activate()
	Deactivate()
	SendToMe(msg Message)
	GetProcessor() Processor
	Drop()
}

type Actor struct {
	address   *PID
	inbox     chan Message
	isClosed  bool
	processor Processor
}

func NewActor(address *PID, processor Processor) Actor {
	return Actor{
		address:   address,
		processor: processor,
		isClosed:  true,
	}
}

func (this *Actor) Activate() {
	if this.IsClosed() {
		this.isClosed = false
		this.inbox = make(chan Message, 10)
		p := this.GetProcessor()
		if p != nil {
			go p.Process(this.inbox)
		}
	}
}

func (this *Actor) GetAddress() *PID {
	return this.address
}

func (this *Actor) IsClosed() bool {
	return this.isClosed
}

func (this *Actor) GetProcessor() Processor {
	return this.processor
}

func (this *Actor) Deactivate() {
	this.isClosed = true
	close(this.inbox)
}

func (this *Actor) SendToMe(msg Message) {
	if !this.isClosed {
		this.inbox <- msg
	} else {
		slog.Warn("lost message for actor", slog.String("to", msg.To.String()), slog.String("from", msg.From.String()))
	}
}

func (this *Actor) Drop() {
	this.Deactivate()
	this.address = nil
	this.processor = nil
	this.inbox = nil
}
