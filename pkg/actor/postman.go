package actor

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Postman struct {
	actors     map[string]*Actor
	context    context.Context
	cancelFunc func()
}

var instance Postman
var onceGuard sync.Once

func GetPostman() *Postman {
	onceGuard.Do(func() {
		ctx, cancFunc := context.WithCancel(context.Background())
		extCancel := make(chan os.Signal, 1)
		signal.Notify(extCancel, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				s := <-extCancel
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					Shutdown()
				}
			}
		}()

		instance = Postman{
			actors:     make(map[string]*Actor, 10),
			context:    ctx,
			cancelFunc: cancFunc,
		}
	})
	return &instance
}

func (this *Postman) GetContext() context.Context {
	return this.context
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

func Shutdown() {
	p := GetPostman()
	for _, a := range p.actors {
		a.Drop()
	}
	p.actors = make(map[string]*Actor)
	p.cancelFunc()
}

func NumActors() int {
	p := GetPostman()
	return len(p.actors)
}
