package services

import (
	"context"
	"log"
	"tdproxy/db"
	"tdproxy/protos"
)

type TradeService struct {
	protos.UnimplementedTradeServiceServer
	TradeDB *db.TradeDB
}

func (s *TradeService) GetTrades(ctx context.Context, request *protos.GetTradesRequest) (resp *protos.GetTradesResponse, err error) {
	trades, err := s.TradeDB.GetTrades(request.TradeIds)
	resp = &protos.GetTradesResponse{
		Trades: make(map[string]*protos.Trade),
	}
	for trade_id, trade := range trades {
		resp.Trades[trade_id] = TradeToProto(trade)
	}
	return resp, nil
}

func (s *TradeService) SaveTrades(ctx context.Context, request *protos.SaveTradesRequest) (resp *protos.SaveTradesResponse, err error) {
	for _, trade := range request.Trades {
		if request.LogTrades {
			log.Println("Saving trade: ", trade.TradeId)
		}
		err = s.TradeDB.SaveTrade(TradeFromProto(trade))
	}
	return &protos.SaveTradesResponse{}, err
}

func (s *TradeService) RemoveTrades(ctx context.Context, request *protos.RemoveTradesRequest) (resp *protos.RemoveTradesResponse, err error) {
	trades, err := s.TradeDB.FilterTrades(
		request.FilterBy.BySymbols,
		request.FilterBy.StartDate,
		request.FilterBy.EndDate,
		request.FilterBy.ByStrategy,
		request.FilterBy.MinGain,
		request.FilterBy.MinProfit,
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	for _, trade := range trades {
		if request.LogTrades {
			log.Println("Removing trade: ", trade.TradeId)
		}
		err = s.TradeDB.RemoveTrade(trade.TradeId)
	}
	resp = &protos.RemoveTradesResponse{}
	return
}

func (s *TradeService) ListTrades(ctx context.Context, request *protos.ListTradesRequest) (resp *protos.ListTradesResponse, err error) {
	resp = &protos.ListTradesResponse{}
	return
}
