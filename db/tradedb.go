package db

import (
	// "log"
	// "errors"
	badger "github.com/dgraph-io/badger/v3"
	"tdproxy/models"
)

type TradeDB struct {
	db *badger.DB
}

func NewTradeDB(db *badger.DB) *TradeDB {
	return &TradeDB{
		db: db,
	}
}

func (authdb *TradeDB) SaveTrade(trade *models.Trade) (err error) {
	return
}
