package services

import (
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/td"
)

type SubscriptionId = int64

type Subscription struct {
	Name    string
	Socket  *td.Socket
	Started bool
}

func NewSubscription(name string, socket *td.Socket) *Subscription {
	return &Subscription{
		Name:   name,
		Socket: socket,
	}
}

func (sub *Subscription) ToProto() *protos.Subscription {
	return &protos.Subscription{
		Name: sub.Name,
	}
}
