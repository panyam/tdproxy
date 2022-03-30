package services

import (
	"context"
	"errors"
	"github.com/panyam/goutils/utils"
	"log"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type ChainService struct {
	protos.UnimplementedChainServiceServer
	TDClient  *tdclient.Client
	AuthStore *tdclient.AuthStore
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
	if !s.AuthStore.EnsureAuthenticated(s.AuthStore.LastAuth().ClientId) {
		return nil, errors.New("Not authenticated.  Call StartLogin first")
	}
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
		}
		for _, option := range chain.Options {
			opt, err := OptionToProto(option)
			if err != nil {
				panic(err)
			}
			if opt != nil {
				resp.Chain.Options = append(resp.Chain.Options, opt)
			}
		}
	}
	return resp, err
}
