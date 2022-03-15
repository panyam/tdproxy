package services

import (
	"context"
	"github.com/panyam/goutils/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/tdclient"
	"log"
)

type AuthService struct {
	protos.UnimplementedAuthServiceServer
	TDClient  *tdclient.Client
	AuthStore *tdclient.AuthStore
}

func (s *AuthService) AddAuthToken(ctx context.Context, request *protos.AddAuthTokenRequest) (*protos.AddAuthTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Execute() not implemented yet")
}

func (s *AuthService) StartLogin(ctx context.Context, request *protos.StartAuthRequest) (*protos.StartAuthResponse, error) {
	s.TDClient.Auth = s.AuthStore.EnsureAuth(request.ClientId, request.CallbackUrl)
	url := s.TDClient.Auth.StartAuthUrl()
	log.Println("StartAuthUrl: ", url)
	if request.LaunchUrl != nil && *request.LaunchUrl {
		utils.OpenBrowser(url)
	}
	return &protos.StartAuthResponse{ContinueAuthUrl: url}, nil
}

func (s *AuthService) CompleteLogin(ctx context.Context, request *protos.CompleteAuthRequest) (*protos.CompleteAuthResponse, error) {
	log.Println("ClientId, CUrl: ", s.TDClient.Auth.ClientId, s.TDClient.Auth.CallbackUrl)
	err := s.TDClient.Auth.CompleteAuth(request.Code)
	if err != nil {
		log.Println("Error completing auth: ", err)
		return nil, err
	}
	s.AuthStore.SaveTokens()
	return &protos.CompleteAuthResponse{Status: true}, nil
}
