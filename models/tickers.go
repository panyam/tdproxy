package models

import (
	"time"
)

type Ticker struct {
	Symbol          string `gorm:"primaryKey"`
	LastRefreshedAt time.Time
	Info            JsonField `gorm:"type:text"`
}

func NewTicker(symbol string, refreshed_at time.Time, info map[string]interface{}) *Ticker {
	return &Ticker{
		Symbol:          symbol,
		LastRefreshedAt: refreshed_at,
		Info:            info,
	}
}
