package services

import (
	"context"
	"errors"
	"github.com/panyam/goutils/utils"
	"log"
	"tdproxy/models"
	"tdproxy/protos"
	"tdproxy/tdclient"
)

type ChainService struct {
	protos.UnimplementedChainServiceServer
	TDClient  *tdclient.Client
	AuthStore *tdclient.AuthStore
}

func (s *ChainService) GetChainInfo(req *protos.GetChainInfoRequest, stream protos.ChainService_GetChainInfoServer) error {
	refresh_type := int32(0)
	if req.RefreshType != nil {
		refresh_type = *req.RefreshType
	}

	for _, symbol := range req.Symbols {
		resp := &protos.GetChainInfoResponse{Symbol: symbol}
		info, err := s.TDClient.GetChainInfo(symbol, refresh_type)
		log.Println("ChainInfo: ", info)
		if err != nil {
			resp.ErrorCode = 1
			resp.ErrorMessage = err.Error()
		} else if info == nil {
			resp.ErrorCode = 2
			resp.ErrorMessage = "Chain not found"
		} else {
			resp.Dates = info.AvailableDates
			resp.LastRefreshedAt = utils.FormatTime(info.LastRefreshedAt)
		}
		if err = stream.Send(resp); err != nil {
			log.Printf("%v.Send(%v) = %v", stream, resp, err)
			return err
		}
	}
	return nil
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
			OptionsByPrice:  make(map[string]*protos.Option),
		}
		for _, option := range chain.Options {
			opt, err := OptionToProto(option)
			if err != nil {
				panic(err)
			}
			if opt != nil {
				resp.Chain.Options = append(resp.Chain.Options, opt)
				resp.Chain.OptionsByPrice[opt.PriceString] = opt
			}
		}

		// Calculate prob range
		if resp.Chain.IsCall {
			resp.Chain.ProbDist = ProbDistToProto(models.DistFromCalls(chain.Options))
		}
	}
	return resp, err
}
