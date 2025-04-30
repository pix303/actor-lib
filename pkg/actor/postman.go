package actor

import (
	"log/slog"
	"sync"
)

type Postman struct {
	actors map[string]*Actor
}

var instance Postman
var onceGuard sync.Once

func GetPostman() *Postman {
	onceGuard.Do(func() {
		instance = Postman{
			actors: make(map[string]*Actor, 10),
		}
	})

	return &instance
}

func RegisterActor(actor *Actor) {
	slog.Info("register an actor", slog.Any("a", actor))
	p := GetPostman()
	p.actors[actor.GetAddress().String()] = actor
}

func DispatchMessage(msg Message) {
	p := GetPostman()
	actor := p.actors[msg.To.String()]
	if actor != nil {
		slog.Debug("actor found, sending msg", slog.String("actor-address", actor.GetAddress().String()), slog.String("to", msg.To.String()))
		actor.Inbox(msg)
	} else {
		slog.Error("actor not found", slog.String("actor-address", msg.To.String()))
	}
}

func DropAllActors() {
	p := GetPostman()
	for _, a := range p.actors {
		a.Drop()
	}
	p.actors = make(map[string]*Actor)
}

func NumActors() int {
	p := GetPostman()
	return len(p.actors)
}
