package services

import (
	"context"
	"encoding/json"
	"fmt"
	pslutils "github.com/panyam/pslite/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/td"
	"legfinder/tdproxy/utils"
	"log"
)

type StreamerService struct {
	protos.UnimplementedStreamerServiceServer
	TDClient     *td.Client
	TopicsFolder string
	topics       map[string]bool
	subs         map[string]*Subscription
	pubsub       *pslutils.PubSub
}

func NewStreamerService(TDClient *td.Client, pubsub *pslutils.PubSub) *StreamerService {
	out := StreamerService{
		TDClient:     TDClient,
		pubsub:       pubsub,
		TopicsFolder: "topics",
		topics:       make(map[string]bool),
		subs:         make(map[string]*Subscription),
	}
	return &out
}

func (s *StreamerService) EnsureTopic(topic_name string) (bool, error) {
	if topic_name == "" {
		return false, nil
	}
	if _, ok := s.topics[topic_name]; ok {
		return true, nil
	}
	// topic doesnt exists so create
	if err := s.pubsub.EnsureTopic(topic_name, s.TopicsFolder); err != nil {
		return false, err
	}
	return true, nil
}

func (s *StreamerService) Subscribe(subreq *protos.SubscribeRequest, stream protos.StreamerService_SubscribeServer) error {
	topic_exists, err := s.EnsureTopic(subreq.TopicName)
	if err != nil {
		return err
	}
	name := subreq.Name
	if sub, ok := s.subs[name]; ok {
		// Already exists and a client is listening to it so return error
		// return status.Error(codes.AlreadyExists, fmt.Sprintf("Subscription (%s) already running", name))
		log.Printf("Subscriptiong %s already connected.  Disconnecting", name)
		sub.Socket.Disconnect()
	}
	newSocket := td.NewSocket(s.TDClient, nil)
	sub := NewSubscription(name, newSocket)
	s.subs[name] = sub
	go sub.Socket.Connect()

	// Now read from the channel that was created and pump it out to our topic
	for {
		newMessage := <-sub.Socket.ReaderChannel()
		if newMessage == nil {
			break
		}
		info, err := structpb.NewStruct(newMessage)
		if err != nil {
			return err
		}
		if topic_exists {
			pubmsg, err := json.Marshal(newMessage)
			if err != nil {
				log.Println("Error encoding json: ", err)
			} else {
				if err := s.pubsub.Publish(subreq.TopicName, pubmsg); err != nil {
					log.Printf("Error publishing message to topic (%s): %v - %v", subreq.TopicName, err, info)
					return err
				}
			}
		}
		msgproto := protos.Message{Info: info}
		if err := stream.Send(&msgproto); err != nil {
			log.Printf("%v.Send(%v) = %v", stream, &msgproto, err)
			return err
		}
	}
	delete(s.subs, name)
	/*
		reply, err := stream.CloseAndRecv()
		if err != nil {
			log.Printf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
			return err
		}
	*/
	return nil
}

func (s *StreamerService) Unsubscribe(ctx context.Context, sub *protos.Subscription) (*protos.EmptyMessage, error) {
	if oldSub, ok := s.subs[sub.Name]; ok {
		if oldSub.Socket.IsRunning() {
			oldSub.Socket.Disconnect()
		}
		delete(s.subs, sub.Name)
	}
	return &protos.EmptyMessage{}, nil
}

func (s *StreamerService) Send(ctx context.Context, request *protos.SendRequest) (*protos.SendResponse, error) {
	log.Println("Received SendRequest: ", request)
	name := request.SubName
	sub, ok := s.subs[name]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Cannot find subscription: %s", name))
	}
	newReq, err := sub.Socket.NewRequest(request.Service, request.Command, true, func(reqparams utils.StringMap) {
		for k, v := range request.Params.AsMap() {
			reqparams[k] = v
		}
	})
	if err != nil {
		return nil, err
	}
	result := sub.Socket.SendRequest(newReq)
	if result {
		return &protos.SendResponse{}, nil
	} else {
		return nil, status.Error(codes.PermissionDenied, "Could not send request")
	}
}
