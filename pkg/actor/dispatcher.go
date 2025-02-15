package actor

import "log/slog"

type Dispatcher interface {
	RegisterActor(actor Dispatchable)
	DispatchMessage(msg Message)
}

type ActorDispatcher struct {
	actors map[string]Dispatchable
}

func NewActorDispatcher() ActorDispatcher {
	return ActorDispatcher{
		actors: make(map[string]Dispatchable, 10),
	}
}

func (this *ActorDispatcher) RegisterActor(actor Dispatchable) (numActors int) {
	this.actors[actor.GetAddress().String()] = actor
	numActors = len(this.actors)
	return
}

func (this *ActorDispatcher) DispatchMessage(msg Message) {
	actor := this.actors[msg.To.String()]
	if actor != nil {
		slog.Debug("find actor and send msg", slog.String("actor-address", actor.GetAddress().String()), slog.String("to", msg.To.String()))
		actor.SendToMe(msg)
	}
}

type Processor interface {
	Process(inbox chan Message)
}
