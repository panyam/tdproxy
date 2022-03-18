package filedb

import (
	"encoding/json"
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"tdproxy/models"
	"time"
)

type TickerDB struct {
	DataRoot          string
	TickersFolderPath string
	tickerCache       map[string]*models.Ticker
}

func NewTickerDB(dataroot string) *TickerDB {
	dataroot, err := filepath.Abs(dataroot)
	if err != nil {
		log.Fatalf("Cannot find directory: %s", dataroot)
	}

	tickers_folder_path := path.Join(dataroot, "tickers")
	if err := os.MkdirAll(tickers_folder_path, 0777); err != nil {
		log.Fatalf("Cannot find directory: %s", dataroot)
	}

	out := TickerDB{DataRoot: dataroot, TickersFolderPath: tickers_folder_path}
	out.tickerCache = make(map[string]*models.Ticker)
	return &out
}

func (db *TickerDB) TickerPathForSymbol(symbol string, ensure bool) (string, error) {
	out := path.Join(db.TickersFolderPath, symbol)
	_, err := os.Stat(out)
	if os.IsNotExist(err) {
		if ensure {
			if err := os.MkdirAll(out, 0777); err != nil {
				log.Fatalf("Cannot create directory: %s", out)
			}
		} else {
			err = fmt.Errorf("Ticker path does not exist")
		}
	}
	out = path.Join(db.TickersFolderPath, symbol, "QUOTES.json")
	_, err = os.Stat(out)
	if os.IsNotExist(err) {
		err = fmt.Errorf("Ticker path does not exist")
	}
	return out, err
}

func (db *TickerDB) GetTicker(symbol string) *models.Ticker {
	ticker_key := symbol
	if val, ok := db.tickerCache[ticker_key]; ok {
		return val
	}

	ticker_path, err := db.TickerPathForSymbol(symbol, false)
	if err != nil {
		return nil
	}

	// Load from file
	data, err := utils.JsonDecodeFile(ticker_path)
	if err != nil {
		return nil
	}

	json_data := data.(map[string]interface{})
	var last_refreshed_at time.Time = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	if datestr, ok := json_data["last_refreshed_at"]; ok {
		last_refreshed_at = utils.ParseTime(datestr.(string))
	}
	ticker_info, ok := json_data["ticker"].(map[string]interface{})
	if !ok {
		return nil
	}
	ticker := models.Ticker{Symbol: symbol,
		LastRefreshedAt: last_refreshed_at,
		Info:            ticker_info}
	db.tickerCache[ticker_key] = &ticker
	return &ticker
}

func (db *TickerDB) SaveTicker(ticker *models.Ticker) error {
	ticker_path, _ := db.TickerPathForSymbol(ticker.Symbol, true)
	content := map[string]interface{}{
		"last_refreshed_at": ticker.LastRefreshedAt,
		"ticker":            ticker.Info,
	}
	marshalled, err := json.Marshal(content)
	if err != nil {
		log.Fatalf("Could not marshall ticker (%s) to JSON", ticker.Symbol)
	}
	d1 := []byte(marshalled)
	if err := os.WriteFile(ticker_path, d1, 0777); err != nil {
		log.Fatalf("Could not write ticker (%s) to file (%s)", ticker.Symbol, ticker_path)
	}
	db.tickerCache[ticker.Symbol] = ticker
	return nil
}
