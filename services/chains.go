package services

import (
	"context"
	"github.com/panyam/goutils/utils"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"log"
	"tdproxy/models"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type ChainService struct {
	protos.UnimplementedChainServiceServer
	TDClient *tdclient.Client
}

func (s *ChainService) GetChainInfo(ctx context.Context, request *protos.GetChainInfoRequest) (*protos.GetChainInfoResponse, error) {
	refresh_type := int32(0)
	if request.RefreshType != nil {
		refresh_type = *request.RefreshType
	}
	info, err := s.TDClient.GetChainInfo(request.Symbol, refresh_type)
	resp := &protos.GetChainInfoResponse{}
	log.Println("Chain: ", info)
	if info != nil {
		resp.Symbol = info.Symbol
		resp.Dates = info.AvailableDates
		resp.LastRefreshedAt = utils.FormatTime(info.LastRefreshedAt)
	}
	return resp, err
}

func (s *ChainService) GetChain(ctx context.Context, request *protos.GetChainRequest) (*protos.GetChainResponse, error) {
	resp := &protos.GetChainResponse{}
	refresh_type := int32(0)
	if request.RefreshType != nil {
		refresh_type = *request.RefreshType
	}
	chain, err := s.TDClient.GetChain(request.Symbol, request.Date, request.IsCall, refresh_type)
	if err == nil {
		resp.Chain = &protos.Chain{
			Symbol:          request.Symbol,
			Date:            request.Date,
			IsCall:          request.IsCall,
			LastRefreshedAt: utils.FormatTime(chain.LastRefreshedAt),
			Options:         make([]*protos.Option, len(chain.Options)),
		}
		for i, option := range chain.Options {
			opt, err := OptionToProto(option)
			if err != nil {
				panic(err)
			}
			resp.Chain.Options[i] = opt
		}
	}
	return resp, err
}

func OptionToProto(option *models.Option) (*protos.Option, error) {
	val, err := option.Info.Value()
	if err != nil {
		return nil, err
	}
	info, err := structpb.NewStruct(val.(utils.StringMap))
	if err != nil {
		return nil, err
	}
	return &protos.Option{
		Symbol:       option.Symbol,
		DateString:   option.DateString,
		PriceString:  option.PriceString,
		IsCall:       option.IsCall,
		AskPrice:     option.AskPrice,
		BidPrice:     option.BidPrice,
		MarkPrice:    option.MarkPrice,
		Multiplier:   option.Multiplier,
		Delta:        option.Delta,
		OpenInterest: option.OpenInterest,
		Info:         info,
	}, nil
}
