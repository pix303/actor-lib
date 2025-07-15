package actor

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ActorNotFoundErr = errors.New("actor not found")
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
	if actor == nil {
		slog.Error("nil actor cant be register")
		return
	}

	slog.Info("register an actor", slog.String("a", actor.GetAddress().String()))
	p := GetPostman()
	p.actors[actor.GetAddress().String()] = actor
	actor.Activate()
}

func SendMessage(msg Message) error {
	p := GetPostman()
	actor := p.actors[msg.To.String()]

	if actor != nil {
		slog.Debug("actor found, sending msg", slog.String("actor-address", msg.To.String()))
		err := actor.Inbox(msg)
		if err != nil {
			slog.Error("actor inbox error", slog.String("actor-address", msg.To.String()), slog.String("error", err.Error()))
			return err
		}
		return nil
	} else {
		slog.Error("actor not found", slog.String("actor-address", msg.To.String()))
		return ActorNotFoundErr
	}
}

func BroadcastMessage(msg Message) {
	p := GetPostman()
	for _, a := range p.actors {
		a.Inbox(msg)
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
