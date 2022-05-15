package db

import (
	"encoding/json"
	"github.com/panyam/goutils/utils"
	"log"
	"sync"
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

func (tdb *TradeDB) GetTradeById(id string) (trade *models.Trade, err error) {
	txn := tdb.db.NewTransaction(true)
	defer txn.Discard()
	item, err := txn.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		tradejson, err := utils.JsonDecodeBytes(val)
		if err != nil {
			return err
		}
		trade = &models.Trade{}
		if err := trade.FromJson(tradejson.(utils.StringMap)); err != nil {
			return err
		}
		return nil
	})
	return
}

func (tdb *TradeDB) GetTrades(trade_ids []string) (trades map[string]*models.Trade, err error) {
	// Now also load their data
	var wg sync.WaitGroup
	trades = make(map[string]*models.Trade)
	for _, trade_id := range trade_ids {
		wg.Add(1)
		go func(trade_id string) {
			defer wg.Done()
			trade, err := tdb.GetTradeById(trade_id)
			if err == nil {
				trades[trade_id] = trade
			}
		}(trade_id)
	}
	wg.Wait()
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
		return nil
	})
	return
}

func (tdb *TradeDB) FilterTrades(
	filter_by_symbols []string,
	filter_start_date string,
	filter_end_date string,
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
	if filter_start_date != "" {
		query = query.Where("date >= ?", filter_start_date)
		// query = query.Where("date IN ?", filter_by_date)
	}
	if filter_end_date != "" {
		query = query.Where("date <= ?", filter_end_date)
	}
	for _, ord := range order_by {
		query = query.Order(ord)
	}
	result := query.Find(&trades)
	err = result.Error
	if err != nil {
		log.Println("Error Filtering Trades: ", err)
	}
	return
}
