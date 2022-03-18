package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/panyam/goutils/utils"
	pslcli "github.com/panyam/pslite/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"log"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type StreamerService struct {
	protos.UnimplementedStreamerServiceServer
	TDClient     *tdclient.Client
	TopicsFolder string
	topics       map[string]bool
	subs         map[string]*Subscription
	pubsub       *pslcli.PubSub
}

func NewStreamerService(TDClient *tdclient.Client, pubsub *pslcli.PubSub) *StreamerService {
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

func (s *StreamerService) Subscribe(subreq *protos.SubscribeRequest, stream protos.StreamerService_SubscribeServer) (err error) {
	topic_exists, err := s.EnsureTopic(subreq.TopicName)
	if err != nil {
		return err
	}
	name := subreq.Name
	if sub, ok := s.subs[name]; ok {
		// Already exists and a client is listening to it so return error
		// return status.Error(codes.AlreadyExists, fmt.Sprintf("Subscription (%s) already running", name))
		sub.Disconnect()
	}
	sub := NewSubscription(name, tdclient.NewSocket(s.TDClient, nil))
	s.subs[name] = sub
	defer delete(s.subs, name)
	if err = sub.StartConnection(); err != nil {
		return err
	}

	// Now read from the channel that was created and pump it out to our topic
	closed := false
	ok := true
	for !closed {
		select {
		case <-stream.Context().Done():
			closed = true
			break
		case newMessage := <-sub.Socket.ReaderChannel():
			if newMessage == nil {
				break
			}
			if topic_exists {
				if ok, err = s.publishToTopic(newMessage, subreq.TopicName); !ok || err != nil {
					break
				}
			}
			if ok, err = s.sendToClient(newMessage, stream); !ok || err != nil {
				break
			}
			break
		}
	}
	sub.Wait()
	return
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

func (s *StreamerService) Send(ctx context.Context, request *protos.SendRequest) (resp *protos.SendResponse, err error) {
	log.Println("Received SendRequest: ", request)
	defer log.Println("Processed SendRequest, err: ", err)
	name := request.SubName
	sub, ok := s.subs[name]
	if !ok {
		err = status.Error(codes.NotFound, fmt.Sprintf("Cannot find subscription: %s", name))
		return
	}
	var newReq utils.StringMap
	newReq, err = sub.Socket.NewRequest(request.Service, request.Command, true, func(reqparams utils.StringMap) {
		for k, v := range request.Params.AsMap() {
			reqparams[k] = v
		}
	})
	if err != nil {
		return
	}
	result := sub.Socket.SendRequest(newReq)
	if result {
		resp = &protos.SendResponse{}
	} else {
		err = status.Error(codes.PermissionDenied, "Could not send request")
	}
	return
}

func (s *StreamerService) publishToTopic(newMessage utils.StringMap, topic_name string) (bool, error) {
	pubmsg, err := json.Marshal(newMessage)
	if err != nil {
		log.Println("Error encoding json: ", err)
	} else {
		err = s.pubsub.Publish(topic_name, pubmsg)
		if err != nil {
			log.Printf("Error publishing message to topic (%s): %v - %v", topic_name, err, pubmsg)
		}
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StreamerService) sendToClient(newMessage utils.StringMap, stream protos.StreamerService_SubscribeServer) (bool, error) {
	info, err := structpb.NewStruct(newMessage)
	if err != nil {
		return false, err
	}
	msgproto := protos.Message{Info: info}
	if err := stream.Send(&msgproto); err != nil {
		log.Printf("%v.Send(%v) = %v", stream, &msgproto, err)
		return false, err
	}
	return true, err
}
