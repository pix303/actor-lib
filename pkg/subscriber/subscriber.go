package subscriber

import "github.com/pix303/actor-lib/pkg/actor"

type Subscribable interface {
	AddSubscription(subscriberAddress *actor.Address)
	NotifySubscribers(msg actor.Message)
}

type SubscriptionsState struct {
	subscribers []*actor.Address
}

func NewSubscribeState() *SubscriptionsState {
	return &SubscriptionsState{
		subscribers: make([]*actor.Address, 0),
	}
}

func (this *SubscriptionsState) AddSubscription(subscriberAddress *actor.Address) {
	this.subscribers = append(this.subscribers, subscriberAddress)
}

func (this *SubscriptionsState) RemoveSubscription(subscriberAddress *actor.Address) {
	for i, v := range this.subscribers {
		if v.IsEqual(subscriberAddress) {
			this.subscribers = append(this.subscribers[:i], this.subscribers[i+1:]...)
		}
	}
}

func (this *SubscriptionsState) NotifySubscribers(msg actor.Message) {
	for _, sub := range this.subscribers {
		msg.To = sub
		actor.SendMessage(msg)
	}
}
