package tdclient

import (
	"errors"
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"tdproxy/db"
	"tdproxy/models"
	"time"
)

var NotAuthenticated = errors.New("Please call StartLogin first")
var InvalidChainJson = errors.New("Invalid chain json dictionary")

type Client struct {
	RootDir              string
	chain_db             db.ChainDB
	ticker_db            db.TickerDB
	DefaultTickerRefresh int32
	DefaultChainRefresh  int32
	Auth                 *Auth
}

const (
	TDAMT_OPT_CHAIN_URL  = "https://api.tdameritrade.com/v1/marketdata/chains"
	TDAMT_OPT_TICKER_URL = "https://api.tdameritrade.com/v1/marketdata/quotes"
)

func NewClient(rootdir string, chain_db db.ChainDB, ticker_db db.TickerDB) *Client {
	fmt.Println("Client Root Dir: ", rootdir)
	out := &Client{RootDir: rootdir}
	os.MkdirAll(rootdir, 0777)
	out.chain_db = chain_db
	out.ticker_db = ticker_db
	out.DefaultTickerRefresh = 1800
	out.DefaultChainRefresh = 1800
	return out
}

func (td *Client) GetTickers(symbols []string, refresh_type int32) (map[string]*models.Ticker, error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultTickerRefresh
	}
	var outdated []string
	tickers := make(map[string]*models.Ticker)
	now := time.Now().UTC()
	for _, sym := range symbols {
		ticker, _ := td.ticker_db.GetTicker(sym)
		if ticker == nil {
			outdated = append(outdated, sym)
		} else if utils.NeedsRefresh(refresh_type, ticker.LastRefreshedAt, now) {
			outdated = append(outdated, sym)
		} else {
			tickers[sym] = ticker
		}
	}

	fetched, err := td.FetchTickers(outdated)
	if err != nil {
		log.Println("Error fetching tickers: ", err)
	} else {
		now2 := time.Now().UTC()
		for sym, ticker := range fetched {
			log.Println("Fetched Ticker: ", sym, ticker)
			tickers[sym] = ticker
			if utils.NeedsRefresh(refresh_type, ticker.LastRefreshedAt, now2) {
				log.Printf("Refresh time not properly set for: %s, Now: %s, LR: %s, RT: %d", sym, utils.FormatTime(now2), utils.FormatTime(ticker.LastRefreshedAt), refresh_type)
			}
		}
	}
	return tickers, err
}

func (td *Client) GetChainInfo(symbol string, refresh_type int32) (*models.ChainInfo, error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultChainRefresh
	}
	chain_info, err := td.chain_db.GetChainInfo(symbol)
	if err != nil {
		log.Println("Error getting chain info: ", symbol, err)
		return nil, err
	}
	now := time.Now().UTC()
	if chain_info == nil || utils.NeedsRefresh(refresh_type, chain_info.LastRefreshedAt, now) {
		if chain_info == nil {
			log.Println("ChainInfo does not exist: ", symbol)
		} else {
			log.Printf("ChainInfo (%s) needs refresh, LastRefreshed: %s, Now: %s, RefreshType: %d",
				symbol, utils.FormatTime(chain_info.LastRefreshedAt), utils.FormatTime(now), refresh_type)
		}
		err = td.FetchChain(symbol, "", true)
		if err == nil {
			// Get the chain info again
			chain_info, err = td.chain_db.GetChainInfo(symbol)
		}
	}
	return chain_info, err
}

func (td *Client) GetChain(symbol string, date string, is_call bool, refresh_type int32) (chain *models.Chain, err error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultChainRefresh
	}
	date = strings.Replace(date, "-", "_", -1)
	now := time.Now().UTC()
	chain, err = td.chain_db.GetChain(symbol, date, is_call)
	if chain == nil || utils.NeedsRefresh(refresh_type, chain.LastRefreshedAt, now) {
		if chain == nil {
			log.Println("Chain does not exist: ", symbol, date, is_call)
		} else {
			log.Printf("Chains (%s-%s-Call(%t)) needs refresh, LastRefreshed: %s, Now: %s, RefreshType: %d",
				symbol, date, is_call, utils.FormatTime(chain.LastRefreshedAt), utils.FormatTime(now), refresh_type)
		}
		err = td.FetchChain(symbol, date, is_call)
		if err == nil {
			chain, err = td.chain_db.GetChain(symbol, date, is_call)
			if chain == nil {
				chtype := "PUT"
				if is_call {
					chtype = "CALL"
				}
				err = fmt.Errorf("No (%s) chain found for %s on %s", chtype, symbol, date)
			} else {
				now2 := time.Now().UTC()
				if utils.NeedsRefresh(refresh_type, chain.LastRefreshedAt, now2) {
					log.Printf("Refresh time not properly set for: %s, Now: %s, LR: %s, RT: %d", symbol, utils.FormatTime(now2), utils.FormatTime(chain.LastRefreshedAt), refresh_type)
				}
			}
		}
		if err != nil {
			log.Println("Error Fetching Chain: ", err)
		}
	}
	return chain, err
}

func (td *Client) FetchTickers(symbols []string) (map[string]*models.Ticker, error) {
	if td.Auth == nil || !td.Auth.IsAuthenticated() {
		return nil, NotAuthenticated
	}
	log.Println("Fetching Tickers from server: ", symbols)
	tickers := make(map[string]*models.Ticker)
	var tail []string
	for len(symbols) > 0 {
		numsyms := 200
		if len(symbols) < numsyms {
			numsyms = len(symbols)
		}
		symbols, tail = symbols[:numsyms], symbols[numsyms:]
		symstr := strings.Join(symbols, "%2C")
		log.Println("Auth: ", td.Auth)
		url := fmt.Sprintf("%s?apikey=%s&symbol=%s", TDAMT_OPT_TICKER_URL, td.Auth.ClientId, symstr)

		result, _, err := utils.JsonGet(url, func(req *http.Request) {
			req.Header.Set("Authorization", td.Auth.Bearer())
		})
		if err != nil {
			return tickers, err
		}
		now := time.Now().UTC()
		for sym, ticker_info := range result.(utils.StringMap) {
			t := models.NewTicker(
				sym,
				now,
				ticker_info.(utils.StringMap),
			)
			tickers[sym] = t
			td.ticker_db.SaveTicker(t)
		}
		symbols = tail
	}
	return tickers, nil
}

func (td *Client) FetchChain(symbol string, date string, is_call bool) error {
	if td.Auth == nil || !td.Auth.IsAuthenticated() {
		return NotAuthenticated
	}
	log.Printf("Fetching chain from server: %s-%s-(Call:%t)", symbol, date, is_call)
	url := fmt.Sprintf("%s?apikey=%s&symbol=%s", TDAMT_OPT_CHAIN_URL, td.Auth.ClientId, symbol)
	// log.Printf("Loading chain data from server for %s: ", url)
	// log.Printf("Bearer Auth: %s", td.Auth.Bearer())
	result, _, err := utils.JsonGet(url, func(req *http.Request) {
		req.Header.Set("Authorization", td.Auth.Bearer())
	})
	if err != nil {
		return err
	}
	calls, puts, err := group_chains_by_date(result.(utils.StringMap), time.Now().UTC())
	if err != nil {
		log.Println("Json Error: ", result)
		return err
	}
	for _, entry := range puts {
		log.Printf("Saving Put Chain for (%s - %s): ", symbol, entry.DateString)
		td.chain_db.SaveChain(entry)
	}
	for _, entry := range calls {
		log.Printf("Saving Call Chain for (%s - %s): ", symbol, entry.DateString)
		td.chain_db.SaveChain(entry)
	}
	// chains = group_chains_by_date(chain)
	return err
}

func extract_chains_by_date(chains_by_date utils.StringMap, symbol string, is_call bool, refreshed_at time.Time) map[string]*models.Chain {
	chains := make(map[string]*models.Chain)
	for date_key, options_by_price := range chains_by_date {
		date := utils.FormatDate(utils.ParseDate(strings.Split(date_key, ":")[0]))
		options := models.ChainFromDict(symbol, date, is_call,
			options_by_price.(utils.StringMap), refreshed_at)
		chains[date] = options
	}
	return chains
}

func group_chains_by_date(chain_json utils.StringMap, refreshed_at time.Time) (
	map[string]*models.Chain,
	map[string]*models.Chain,
	error) {
	if chain_json["error"] != nil {
		return nil, nil, errors.New(chain_json["error"].(string))
	}
	if chain_json["symbol"] == nil {
		return nil, nil, InvalidChainJson
	}
	ticker := chain_json["symbol"].(string)
	put_exp_date_map := chain_json["putExpDateMap"].(utils.StringMap)
	call_exp_date_map := chain_json["callExpDateMap"].(utils.StringMap)
	calls := extract_chains_by_date(call_exp_date_map, ticker, true, refreshed_at)
	puts := extract_chains_by_date(put_exp_date_map, ticker, false, refreshed_at)
	return calls, puts, nil
}
