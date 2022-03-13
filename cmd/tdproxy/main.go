package main

import (
	"flag"
	"fmt"
	"github.com/panyam/pslite/cli"
	pslutils "github.com/panyam/pslite/utils"
	"google.golang.org/grpc"
	"legfinder/tdproxy/protos"
	svc "legfinder/tdproxy/services"
	"legfinder/tdproxy/tdclient"
	"legfinder/tdproxy/utils"
	"log"
	"net"
	"strconv"
	"strings"
)

const TEST_CALLBACK_URL = "https://localhost:8000/callback"
const TEST_CLIENT_ID = ""

var (
	port           = flag.Int("port", utils.DefaultServerPort(), "Port on which gRPC server should listen TCP conn.")
	tdroot         = flag.String("tdroot", "~/.tdroot", "Root location of where TD data is downloaded too")
	client_id      = flag.String("client_id", TEST_CLIENT_ID, "TD Ameritrade Client ID")
	callback_port  = flag.Int("callback_port", utils.DefaultCallbackPort(), "Port on which OAuth Callback handler listen on.")
	callback_url   = flag.String("callback_url", TEST_CALLBACK_URL, "TD Ameritrade Auth Callback URl")
	callback_cert  = flag.String("callback_cert", "./tdclient/server.crt", "Certificate file for SSL Callback handler")
	callback_pkey  = flag.String("callback_pkey", "./tdclient/server.key", "Private key file for SSL Callback handler")
	topic_endpoint = flag.String("topic_endpoint", fmt.Sprintf("%d", pslutils.DefaultServerPort()), "End point where topics can be published and subscribed to")
	topics_folder  = flag.String("topics_folder", "~/.tdroot/topics", "End point where topics can be published and subscribed to")
)

func createPubsubClient() *cli.PubSub {
	pubsub, err := cli.NewPubSub(*topic_endpoint)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Index(*topic_endpoint, ":") < 0 {
		// Start us locally on the given port
		topicport, err := strconv.ParseInt(*topic_endpoint, 10, 64)
		if err != nil {
			log.Fatal(fmt.Errorf("Invalid topic_endpoint port: %s", *topic_endpoint))
		}
		log.Printf("Serving pubsub topics on %d - %s", topicport, *topics_folder)
		go cli.PSLServe(int(topicport), *topics_folder)
	}
	if err != nil {
		log.Fatal(err)
	}
	return pubsub
}

func main() {
	flag.Parse()
	grpcServer := grpc.NewServer()
	// see if we need to start the pubsub endpoint locally
	pubsub := createPubsubClient()
	tdinfo := tdclient.NewClient(utils.ExpandUserPath(*tdroot), *client_id, *callback_url)
	callbackHandler := tdclient.NewCallbackHandler(tdinfo,
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
	log.Printf("Initializing TDProxy gRPC server on port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
