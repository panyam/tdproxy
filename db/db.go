package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"legfinder/tdproxy/models"
	"legfinder/tdproxy/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type QuoteDB struct {
	DataRoot          string
	TickersFolderPath string
	chainInfoCache    map[string]*models.TickerChainInfo
	chainCache        map[string]*models.Chain
	tickerCache       map[string]*models.Ticker
}

func NewDB(dataroot string) *QuoteDB {
	dataroot, err := filepath.Abs(dataroot)
	if err != nil {
		log.Fatalf("Cannot find directory: %s", dataroot)
	}

	tickers_folder_path := path.Join(dataroot, "tickers")
	if err := os.MkdirAll(tickers_folder_path, 0777); err != nil {
		log.Fatalf("Cannot find directory: %s", dataroot)
	}

	out := QuoteDB{DataRoot: dataroot, TickersFolderPath: tickers_folder_path}
	out.chainInfoCache = make(map[string]*models.TickerChainInfo)
	out.chainCache = make(map[string]*models.Chain)
	out.tickerCache = make(map[string]*models.Ticker)
	return &out
}

func (db *QuoteDB) TickerPathForSymbol(symbol string, ensure bool) (string, error) {
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

func (db *QuoteDB) ChainInfoPathForSymbol(symbol string, ensure bool) (string, error) {
	out := path.Join(db.TickersFolderPath, symbol, "chains")
	var err error = nil
	if _, err = os.Stat(out); os.IsNotExist(err) {
		if ensure {
			if err = os.MkdirAll(out, 0777); err != nil {
				log.Fatalf("Cannot create directory: %s", out)
			}
		} else {
			err = fmt.Errorf("Chain info path for %s does not exist", symbol)
		}
	}
	out = path.Join(out, "Info.json")
	if err == nil {
		_, err = os.Stat(out)
	}
	log.Println("Info Path for: ", symbol, out)
	return out, err
}

func (db *QuoteDB) ChainPathForSymbol(symbol string, date string, ensure bool) (string, error) {
	out := path.Join(db.TickersFolderPath, symbol, "chains")
	if len(strings.TrimSpace(date)) > 0 {
		out = path.Join(out, date)
	}
	var err error = nil
	if _, err = os.Stat(out); os.IsNotExist(err) {
		if ensure {
			if err = os.MkdirAll(out, 0777); err != nil {
				log.Fatalf("Cannot create directory: %s", out)
			}
		} else {
			err = fmt.Errorf("Chain path for %s on %s does not exist", symbol, date)
		}
	}
	return out, err
}

func chainTypeFor(is_call bool) string {
	chtype := "puts"
	if is_call {
		chtype = "calls"
	}
	return chtype
}

func (db *QuoteDB) ChainKeyFor(symbol string, date string, is_call bool) string {
	chtype := chainTypeFor(is_call)
	return strings.Join([]string{symbol, date, chtype}, "/")
}

/**
 * Saves metadata about a chain - currently when it was last refreshed.
 */
func (db *QuoteDB) SaveChainInfo(symbol string, last_refreshed_at time.Time) error {
	chain_info_path, err := db.ChainInfoPathForSymbol(symbol, true)
	d1 := []byte(fmt.Sprintf(`{"last_refreshed_at": "%s"}`, utils.FormatTime(last_refreshed_at)))
	err = os.WriteFile(chain_info_path, d1, 0777)
	if err != nil {
		log.Fatalf("Could not write chain info (%s) to file (%s)", symbol, chain_info_path)
	}
	return err
}

/**
 * Get information about a chain.
 */
func (db *QuoteDB) GetChainInfo(symbol string) (*models.TickerChainInfo, error) {
	chain_info_path, err := db.ChainInfoPathForSymbol(symbol, true)
	last_refreshed_at := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	if err == nil {
		info, _ := utils.JsonDecodeFile(chain_info_path)
		if info != nil {
			last_refreshed_at = utils.ParseTime(info.(map[string]interface{})["last_refreshed_at"].(string))
		}
	}
	chain_path, err := db.ChainPathForSymbol(symbol, "", true)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(chain_path)
	if err != nil {
		log.Fatal(err)
	}

	out := &models.TickerChainInfo{
		Symbol: symbol,
	}
	for _, file := range files {
		fname := strings.Replace(file.Name(), "_", "-", -1)
		_, err := time.Parse(utils.DATE_FORMAT, fname)
		if err == nil {
			out.AvailableDates = append(out.AvailableDates, fname)
		}
	}
	out.LastRefreshedAt = last_refreshed_at
	return out, nil
}

func (db *QuoteDB) GetChain(symbol string, date string, is_call bool) *models.Chain {
	chain_key := db.ChainKeyFor(symbol, date, is_call)
	if val, ok := db.chainCache[chain_key]; ok {
		return val
	}

	// Get folder when call and put chains exist
	chain_folder, err := db.ChainPathForSymbol(symbol, date, false)
	if err != nil {
		return nil
	}

	chtype := chainTypeFor(is_call)
	chain_path := path.Join(chain_folder, fmt.Sprintf("%s.json", chtype))
	// Load from file
	contents, err := os.ReadFile(chain_path)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(contents)))
	var json_data map[string]interface{}
	if err := decoder.Decode(&json_data); err != nil {
		log.Fatal("Invalid error decoding json: ", err)
		return nil
	}

	var last_refreshed_at time.Time = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	if datestr, ok := json_data["last_refreshed_at"]; ok {
		last_refreshed_at = utils.ParseTime(datestr.(string))
	}
	options_by_price, ok := json_data["chain"].(map[string]interface{})
	if !ok {
		return nil
	}
	chain := models.ChainFromDict(symbol, date, is_call, options_by_price, last_refreshed_at)
	db.chainCache[chain_key] = chain
	return chain
}

func (db *QuoteDB) SaveChain(chain *models.Chain) error {
	chtype := chainTypeFor(chain.IsCall)
	log.Printf("Saving (%s) chain for %s on %s\n", chtype, chain.Symbol, chain.DateString)
	chain_path, err := db.ChainPathForSymbol(chain.Symbol, chain.DateString, true)
	chain_path = path.Join(chain_path, fmt.Sprintf("%s.json", chtype))
	chains := make(map[string]interface{})
	for _, option := range chain.Options {
		chains[option.PriceString] = option.Info
	}
	content := map[string]interface{}{
		"last_refreshed_at": utils.FormatTime(chain.LastRefreshedAt),
		"chain":             chains,
	}
	marshalled, err := json.Marshal(content)
	if err != nil {
		log.Fatalf("Could not marshall chain (%s) to JSON", chain.Symbol)
	}
	d1 := []byte(marshalled)
	if err := os.WriteFile(chain_path, d1, 0777); err != nil {
		log.Fatalf("Could not write chain (%s) to file (%s)", chain.Symbol, chain_path)
	}
	chain_key := db.ChainKeyFor(chain.Symbol, chain.DateString, chain.IsCall)
	db.chainCache[chain_key] = chain
	return nil
}

func (db *QuoteDB) GetTicker(symbol string) *models.Ticker {
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

func (db *QuoteDB) SaveTicker(ticker *models.Ticker) error {
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
