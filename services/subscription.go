package services

import (
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/tdclient"
	"log"
)

type SubscriptionId = int64

type Subscription struct {
	Name    string
	Socket  *tdclient.Socket
	Started bool
}

func NewSubscription(name string, socket *tdclient.Socket) *Subscription {
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

func (sub *Subscription) Disconnect() {
	log.Printf("Disconnecting Subscription: %s", sub.Name)
	sub.Socket.Disconnect()
}

/**
 * Starts a connection on the socket so that messages can be read and written.
 */
func (sub *Subscription) StartConnection() error {
	return sub.Socket.StartConnection()
}

func (sub *Subscription) Wait() {
	sub.Socket.WaitForFinish()
}
