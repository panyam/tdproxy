package services

import (
	"context"
	"errors"
	"log"
	// "fmt"
	"github.com/panyam/goutils/utils"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type TickerService struct {
	protos.UnimplementedTickerServiceServer
	TDClient  *tdclient.Client
	AuthStore *tdclient.AuthStore
}

func (s *TickerService) GetTickers(ctx context.Context, request *protos.GetTickersRequest) (*protos.GetTickersResponse, error) {
	lastAuth := s.AuthStore.LastAuth()
	if lastAuth == nil || !s.AuthStore.EnsureAuthenticated(lastAuth.ClientId) {
		return nil, errors.New("Not authenticated.  Call StartLogin first")
	}
	refresh_type := int32(0)
	if request.RefreshType != nil {
		refresh_type = *request.RefreshType
	}
	log.Println("Fetching Tickers: ", request.Symbols)
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
