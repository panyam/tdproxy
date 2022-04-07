package db

import (
	"encoding/json"
	"log"
	// "errors"
	badger "github.com/dgraph-io/badger/v3"
	"gorm.io/gorm"
	"tdproxy/models"
)

type TradeDB struct {
	db      *badger.DB
	indexdb *gorm.DB
}

func NewTradeDB(db *badger.DB, indexdb *gorm.DB) *TradeDB {
	indexdb.AutoMigrate(&models.Trade{})
	return &TradeDB{
		db:      db,
		indexdb: indexdb,
	}
}

func (tdb *TradeDB) Close() (err error) {
	tdb.db.Close()
	return
}

func (tdb *TradeDB) SaveTrade(trade *models.Trade) (err error) {
	// TODO - Set into KV
	tjson, err := json.Marshal(trade.ToJson())
	if err != nil {
		log.Println("Json Encoding Error: ", err)
		return err
	}
	err = tdb.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(trade.TradeId), []byte(tjson)); err != nil {
			log.Println("Set Err: ", err)
			return err
		}
		return nil
	})

	if err == nil {
		result := tdb.indexdb.Updates(trade)
		err = result.Error
		if err == nil && result.RowsAffected == 0 {
			err = tdb.indexdb.Create(trade).Error
		}
	}
	if err != nil {
		log.Println("Error Saving Trade: ", err)
	}
	return
}

func (tdb *TradeDB) RemoveTrade(tid string) (err error) {
	// remove from index and remove from DB
	var trade models.Trade
	err = tdb.indexdb.Where("trade_id = ?", tid).Delete(&trade).Error

	// TODO - Delete from KV
	err = tdb.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(trade.TradeId)); err != nil {
			return err
		}
		return txn.Commit()
	})
	return
}

func (tdb *TradeDB) FilterTrades(
	filter_by_symbols []string,
	filter_by_date []string,
	filter_by_strategy []string,
	min_gain float64,
	min_profit float64,
	order_by []string,
) (trades []*models.Trade, err error) {
	query := tdb.indexdb.Where("payoff_expected_gain >= ?", min_profit)
	query = query.Where("payoff_gain_prob >= ?", min_gain)
	if len(filter_by_symbols) > 0 {
		query = query.Where("symbol IN ?", filter_by_symbols)
	}
	if len(filter_by_strategy) > 0 {
		query = query.Where("strategy IN ?", filter_by_strategy)
	}
	if len(filter_by_date) > 0 {
		query = query.Where("date IN ?", filter_by_date)
	}
	for _, ord := range order_by {
		query = query.Order(ord)
	}
	result := query.Find(&trades)
	err = result.Error
	return
}
