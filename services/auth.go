package services

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"legfinder/tdproxy/protos"
	"legfinder/tdproxy/tdclient"
	"legfinder/tdproxy/utils"
	"log"
)

type AuthService struct {
	protos.UnimplementedAuthServiceServer
	TDClient *tdclient.Client
}

func (s *AuthService) AddAuthToken(ctx context.Context, request *protos.AddAuthTokenRequest) (*protos.AddAuthTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Execute() not implemented yet")
}

func (s *AuthService) StartLogin(ctx context.Context, request *protos.StartAuthRequest) (*protos.StartAuthResponse, error) {
	s.TDClient.Auth = tdclient.NewAuth(s.TDClient.RootDir, request.ClientId, request.CallbackUrl)
	url := s.TDClient.Auth.StartAuthUrl()
	log.Println("StartAuthUrl: ", url)
	if request.LaunchUrl != nil && *request.LaunchUrl {
		utils.OpenBrowser(url)
	}
	return &protos.StartAuthResponse{ContinueAuthUrl: url}, nil
}

func (s *AuthService) CompleteLogin(ctx context.Context, request *protos.CompleteAuthRequest) (*protos.CompleteAuthResponse, error) {
	log.Println("ClientId, CUrl: ", s.TDClient.Auth.ClientId, s.TDClient.Auth.CallbackUrl)
	_, err := s.TDClient.Auth.CompleteAuth(request.Code)
	if err != nil {
		return nil, err
	}
	return &protos.CompleteAuthResponse{Status: true}, nil
}
