package services

import (
	"context"
	// "fmt"
	"github.com/panyam/goutils/utils"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type TickerService struct {
	protos.UnimplementedTickerServiceServer
	TDClient *tdclient.Client
}

func (s *TickerService) GetTickers(ctx context.Context, request *protos.GetTickersRequest) (*protos.GetTickersResponse, error) {
	refresh_type := int32(0)
	if request.RefreshType != nil {
		refresh_type = *request.RefreshType
	}
	tickers, err := s.TDClient.GetTickers(request.Symbols, refresh_type)
	resp := &protos.GetTickersResponse{
		Errors:  make(map[string]string),
		Tickers: make(map[string]*protos.Ticker),
	}
	for sym, ticker := range tickers {
		info, err := structpb.NewStruct(ticker.Info)
		if err != nil {
			resp.Errors[sym] = err.Error()
		}
		tickerproto := &protos.Ticker{
			Symbol:          sym,
			LastRefreshedAt: utils.FormatTime(ticker.LastRefreshedAt),
			Info:            info,
		}
		resp.Tickers[sym] = tickerproto
	}
	return resp, err
}
