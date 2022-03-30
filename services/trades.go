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
	trades, err := s.TradeDB.FilterTrades(
		request.FilterBy.BySymbols,
		request.FilterBy.ByDate,
		request.FilterBy.ByStrategy,
		request.FilterBy.MinGain,
		request.FilterBy.MinProfit,
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	for _, trade := range trades {
		err = s.TradeDB.RemoveTrade(trade.TradeId)
	}
	return
}

func (s *TradeService) ListTrades(ctx context.Context, request *protos.ListTradesRequest) (resp *protos.ListTradesResponse, err error) {
	return
}
