package main

import (
	"flag"
	"fmt"
	pslutils "github.com/panyam/pslite/utils"
	"google.golang.org/grpc"
	"legfinder/tdproxy/protos"
	svc "legfinder/tdproxy/services"
	"legfinder/tdproxy/td"
	"legfinder/tdproxy/utils"
	"log"
	"net"
)

const TEST_CALLBACK_URL = "https://localhost:8000/callback"
const TEST_CLIENT_ID = ""

var (
	port           = flag.Int("port", utils.DefaultServerPort(), "Port on which gRPC server should listen TCP conn.")
	tdroot         = flag.String("tdroot", "~/.tdroot", "Root location of where TD data is downloaded too")
	client_id      = flag.String("client_id", TEST_CLIENT_ID, "TD Ameritrade Client ID")
	callback_port  = flag.Int("callback_port", utils.DefaultCallbackPort(), "Port on which OAuth Callback handler listen on.")
	callback_url   = flag.String("callback_url", TEST_CALLBACK_URL, "TD Ameritrade Auth Callback URl")
	callback_cert  = flag.String("callback_cert", "./td/server.crt", "Certificate file for SSL Callback handler")
	callback_pkey  = flag.String("callback_pkey", "./td/server.key", "Private key file for SSL Callback handler")
	topic_endpoint = flag.String("topic_endpoint", pslutils.DefaultServerAddress(), "End point where topics can be published and subscribed to")
	topics_folder  = flag.String("topics_folder", "~/.tdroot/topics", "End point where topics can be published and subscribed to")
)

func main() {
	flag.Parse()
	grpcServer := grpc.NewServer()
	pubsub, err := pslutils.NewPubSub(*topic_endpoint)
	if err != nil {
		log.Fatal(err)
	}
	tdinfo := td.NewClient(utils.ExpandUserPath(*tdroot), *client_id, *callback_url)
	callbackHandler := td.NewCallbackHandler(tdinfo,
		*callback_port,
		*callback_cert,
		*callback_pkey)
	go callbackHandler.Start()
	protos.RegisterTickerServiceServer(grpcServer, &svc.TickerService{TDClient: tdinfo})
	protos.RegisterAuthServiceServer(grpcServer, &svc.AuthService{TDClient: tdinfo})
	protos.RegisterChainServiceServer(grpcServer, &svc.ChainService{TDClient: tdinfo})

	streamer_svc := svc.NewStreamerService(tdinfo, pubsub)
	streamer_svc.TopicsFolder = *topics_folder
	protos.RegisterStreamerServiceServer(grpcServer, streamer_svc)
	log.Printf("Initializing gRPC server on port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
