package models

import (
	"time"
)

type TickerJsonField struct {
	*Json
	TickerSymbol string
}

type Ticker struct {
	Symbol          string `gorm:"primaryKey"`
	LastRefreshedAt time.Time
	Info            TickerJsonField // *Json
}

func NewTicker(symbol string, refreshed_at time.Time, info map[string]interface{}) *Ticker {
	return &Ticker{
		Symbol:          symbol,
		LastRefreshedAt: refreshed_at,
		Info: TickerJsonField{
			TickerSymbol: symbol,
			Json:         NewJson(info),
		},
	}
}
