package services

import (
	"context"
	"github.com/panyam/goutils/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"tdproxy/protos"
	"tdproxy/tdclient"
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
	auth, err := s.AuthStore.EnsureAuth(request.ClientId, request.CallbackUrl)
	if err != nil {
		log.Println("Login Error: ", err)
		return nil, err
	}
	s.TDClient.Auth = auth
	log.Println("After StartAuth: ", s.TDClient.Auth)
	url := s.TDClient.Auth.StartAuthUrl()
	log.Println("StartAuthUrl: ", url)
	if request.LaunchUrl != nil && *request.LaunchUrl {
		utils.OpenBrowser(url)
	}
	return &protos.StartAuthResponse{ContinueAuthUrl: url}, nil
}

func (s *AuthService) CompleteLogin(ctx context.Context, request *protos.CompleteAuthRequest) (*protos.CompleteAuthResponse, error) {
	log.Println("ClientId, CUrl: ", s.TDClient.Auth.ClientId, s.TDClient.Auth.CallbackUrl, s.TDClient.Auth.ExpiresAt)
	err := s.TDClient.Auth.CompleteAuth(request.Code)
	if err != nil {
		log.Println("Error completing auth: ", err)
		return nil, err
	}
	log.Println("Completed Auth: ", s.TDClient.Auth)
	s.AuthStore.SaveAuth(s.TDClient.Auth)
	return &protos.CompleteAuthResponse{Status: true}, nil
}
