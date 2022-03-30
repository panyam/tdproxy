package services

import (
	"context"
	"tdproxy/db"
	"tdproxy/protos"
)

type TradeService struct {
	protos.UnimplementedTradeServiceServer
	TradeDB *db.TradeDB
}

func (s *TradeService) SaveTrades(ctx context.Context, request *protos.SaveTradesRequest) (resp *protos.SaveTradesResponse, err error) {
	for _, trade := range request.Trades {
		err = s.TradeDB.SaveTrade(TradeFromProto(trade))
	}
	return
}

func (s *TradeService) RemoveTrades(ctx context.Context, request *protos.RemoveTradesRequest) (resp *protos.RemoveTradesResponse, err error) {
	keys := s.TradeDB.FilterTrades(request.FilterBy)
	for _, tid := range keys {
		err = s.TradeDB.RemoveTrade(tid)
	}
	return
}

func (s *TradeService) ListTrades(ctx context.Context, request *protos.ListTradesRequest) (resp *protos.ListTradesResponse, err error) {
	return
}
