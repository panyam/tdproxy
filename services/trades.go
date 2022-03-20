package services

import (
	// "context"
	// "github.com/panyam/goutils/utils"
	// structpb "google.golang.org/protobuf/types/known/structpb"
	// "log"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type TradeService struct {
	protos.UnimplementedTradeServiceServer
	TDClient *tdclient.Client
}
