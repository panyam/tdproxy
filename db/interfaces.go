package db

import (
	"tdproxy/models"
	"time"
)

type AuthDB interface {
	EnsureAuth(client_id string) (*models.Auth, error)
	SaveAuth(auth *models.Auth) error
	LastAuth() *models.Auth
}

type ChainDB interface {
	GetChainInfo(symbol string) (*models.ChainInfo, error)
	SaveChainInfo(symbol string, last_refreshed_at time.Time) error
	GetChain(symbol string, date string, is_call bool) (*models.Chain, error)
	SaveChain(chain *models.Chain) error
}

type TickerDB interface {
	GetTicker(symbol string) (*models.Ticker, error)
	SaveTicker(ticker *models.Ticker) error
}
