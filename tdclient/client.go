package tdclient

import (
	"errors"
	"fmt"
	"legfinder/tdproxy/db"
	"legfinder/tdproxy/models"
	"legfinder/tdproxy/utils"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var NotAuthenticated = errors.New("Please call StartLogin first")

type Client struct {
	RootDir              string
	db                   *db.QuoteDB
	DefaultTickerRefresh int32
	DefaultChainRefresh  int32
	Auth                 *Auth
}

const (
	TDAMT_OPT_CHAIN_URL  = "https://api.tdameritrade.com/v1/marketdata/chains"
	TDAMT_OPT_TICKER_URL = "https://api.tdameritrade.com/v1/marketdata/quotes"
)

func NewClient(rootdir string, client_id string, callback_url string) *Client {
	fmt.Println("Client Root Dir: ", rootdir)
	out := &Client{RootDir: rootdir}
	os.MkdirAll(rootdir, 0777)
	out.db = db.NewDB(path.Join(rootdir, "quotedb"))
	out.DefaultTickerRefresh = 1800
	out.DefaultChainRefresh = 1800

	if client_id != "" && callback_url != "" {
		out.Auth = NewAuth(rootdir, client_id, callback_url)
	}
	return out
}

func (td *Client) GetTickers(symbols []string, refresh_type int32) (map[string]*models.Ticker, error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultTickerRefresh
	}
	qdb := td.db
	var outdated []string
	tickers := make(map[string]*models.Ticker)
	now := time.Now().UTC()
	for _, sym := range symbols {
		ticker := qdb.GetTicker(sym)
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
			log.Println("Sym: ", sym, ticker)
			tickers[sym] = ticker
			if utils.NeedsRefresh(refresh_type, ticker.LastRefreshedAt, now2) {
				log.Printf("Refresh time not properly set for: %s, Now: %s, LR: %s, RT: %d", sym, utils.FormatTime(now2), utils.FormatTime(ticker.LastRefreshedAt), refresh_type)
			}
		}
	}
	return tickers, err
}

func (td *Client) GetChainInfo(symbol string, refresh_type int32) (*models.TickerChainInfo, error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultChainRefresh
	}
	chain_info, err := td.db.GetChainInfo(symbol)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if chain_info == nil || utils.NeedsRefresh(refresh_type, chain_info.LastRefreshedAt, now) {
		err = td.FetchChain(symbol, "", true)
		if err == nil {
			err = td.db.SaveChainInfo(symbol, time.Now().UTC())
			// Get the chain info again
			chain_info, err = td.db.GetChainInfo(symbol)
		}
	}
	return chain_info, err
}

func (td *Client) GetChain(symbol string, date string, is_call bool, refresh_type int32) (chain *models.Chain, err error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultChainRefresh
	}
	date = strings.Replace(date, "-", "_", -1)
	qdb := td.db
	now := time.Now().UTC()
	chain = qdb.GetChain(symbol, date, is_call)
	if chain == nil || utils.NeedsRefresh(refresh_type, chain.LastRefreshedAt, now) {
		err = td.FetchChain(symbol, date, is_call)
		if err == nil {
			chain = qdb.GetChain(symbol, date, is_call)
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
	}
	return chain, err
}

func (td *Client) FetchTickers(symbols []string) (map[string]*models.Ticker, error) {
	if td.Auth == nil {
		return nil, NotAuthenticated
	}
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
			tickers[sym] = &models.Ticker{
				Symbol:          sym,
				Info:            ticker_info.(utils.StringMap),
				LastRefreshedAt: now,
			}
		}
		symbols = tail
	}
	return tickers, nil
}

func (td *Client) FetchChain(symbol string, date string, is_call bool) error {
	if td.Auth == nil {
		return NotAuthenticated
	}
	log.Printf("Loading chain data from server for %s...", symbol)
	url := fmt.Sprintf("%s?apikey=%s&symbol=%s", TDAMT_OPT_CHAIN_URL, td.Auth.ClientId, symbol)
	result, _, err := utils.JsonGet(url, func(req *http.Request) {
		req.Header.Set("Authorization", td.Auth.Bearer())
	})
	if err != nil {
		return err
	}
	calls, puts := group_chains_by_date(result.(utils.StringMap), time.Now().UTC())
	for _, entry := range calls {
		td.db.SaveChain(entry)
	}
	for _, entry := range puts {
		td.db.SaveChain(entry)
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
	map[string]*models.Chain) {
	ticker := chain_json["symbol"].(string)
	put_exp_date_map := chain_json["putExpDateMap"].(utils.StringMap)
	call_exp_date_map := chain_json["callExpDateMap"].(utils.StringMap)
	calls := extract_chains_by_date(call_exp_date_map, ticker, true, refreshed_at)
	puts := extract_chains_by_date(put_exp_date_map, ticker, false, refreshed_at)
	return calls, puts
}
