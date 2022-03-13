package services

import (
	"context"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/tdclient"
	"legfinder/tdproxy/utils"
	"log"
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
			info, err := structpb.NewStruct(option.Info)
			if err != nil {
				panic(err)
			}
			resp.Chain.Options[i] = &protos.Option{
				Symbol:      option.Symbol,
				DateString:  option.DateString,
				PriceString: option.PriceString,
				IsCall:      option.IsCall,
				Info:        info,
			}
		}
	}
	return resp, err
}
