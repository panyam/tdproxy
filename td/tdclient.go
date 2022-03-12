package td

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"legfinder/tdproxy/db"
	"legfinder/tdproxy/models"
	"legfinder/tdproxy/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Client struct {
	RootDir              string
	ClientId             string
	CallbackUrl          string
	db                   *db.QuoteDB
	DefaultTickerRefresh int32
	DefaultChainRefresh  int32
	authToken            utils.StringMap
	userPrincipals       utils.StringMap
	credentials          utils.StringMap
	wsUrl                *url.URL
}

const (
	TDAMT_AUTH_URL            = "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=%s&client_id=%s%%40AMER.OAUTHAP"
	TDAMT_TOKEN_URL           = "https://api.tdameritrade.com/v1/oauth2/token"
	TDAMT_OPT_CHAIN_URL       = "https://api.tdameritrade.com/v1/marketdata/chains"
	TDAMT_OPT_TICKER_URL      = "https://api.tdameritrade.com/v1/marketdata/quotes"
	TDAMT_USER_PRINCIPALS_URL = "https://api.tdameritrade.com/v1/userprincipals"
)

func NewClient(rootdir string, client_id string, callback_url string) *Client {
	fmt.Println("Client Root Dir: ", rootdir)
	out := &Client{RootDir: rootdir, ClientId: client_id, CallbackUrl: callback_url}
	os.MkdirAll(rootdir, 0777)
	out.db = db.NewDB(path.Join(rootdir, "quotedb"))
	out.ReloadAuthTokens()
	out.DefaultTickerRefresh = 1800
	out.DefaultChainRefresh = 1800
	return out
}

func (td *Client) IsAuthenticated() bool {
	authToken := td.ReloadAuthTokens()
	if authToken == nil {
		return false
	}
	expires_at_str := authToken["expires_at"]
	if expires_at_str == nil {
		return false
	}
	expires_at := utils.ParseTime(expires_at_str.(string))
	now := time.Now().UTC()
	if now.Sub(expires_at) >= 0 {
		return false
	}
	return true
}

func (td *Client) TokenFilePath() string {
	return path.Join(td.RootDir, "tokens")
}

func (td *Client) Bearer() string {
	return fmt.Sprintf("Bearer %s", td.GetAccessToken())
}

func (td *Client) GetAccessToken() string {
	return td.authToken["access_token"].(string)
}

func (td *Client) ReloadAuthTokens() utils.StringMap {
	td.authToken = nil
	contents, err := os.ReadFile(td.TokenFilePath())
	if err != nil {
		log.Println(err)
		return nil
	}
	token, err := utils.JsonDecodeBytes(contents)
	if err != nil {
		log.Println(err)
		return nil
	}
	td.authToken = token.(utils.StringMap)
	fmt.Println("Loaded auth token: ", token)
	return td.authToken
}

func (td *Client) StartAuthUrl() string {
	return td.StartAuthUrlFor(td.ClientId, td.CallbackUrl)
}

func (td *Client) StartAuthUrlFor(client_id string, callback_url string) string {
	callback_quoted := url.QueryEscape(callback_url)
	url := fmt.Sprintf(TDAMT_AUTH_URL, callback_quoted, client_id)
	return url
}

func (td *Client) CompleteAuth(code string) (utils.StringMap, error) {
	return td.CompleteAuthFor(td.ClientId, td.CallbackUrl, code)
}

func (td *Client) CompleteAuthFor(client_id string,
	callback_url string, code string) (utils.StringMap, error) {
	// decoded := code
	decoded, err := url.PathUnescape(code)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", client_id)
	form.Add("redirect_uri", callback_url)
	form.Add("code", decoded)
	fmt.Println("Form: ", form)
	postreq, err := http.NewRequest("POST", TDAMT_TOKEN_URL, strings.NewReader(form.Encode()))
	postreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(postreq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	if response.StatusCode != 200 {
		fmt.Println("Failed Response: ", response)
		return nil, fmt.Errorf(response.Status)
	}

	// Save it locally now
	token, err := utils.JsonDecodeBytes(body)
	if err != nil {
		fmt.Println("Invalid response json: ", err)
		return nil, err
	}
	tokenmap := token.(utils.StringMap)
	expires_in := time.Duration(tokenmap["expires_in"].(float64))
	expires_at := now.Add(expires_in * time.Second)
	log.Println("Now, ExpiresIn, ExpiresAt: ", now, expires_in, expires_at)
	tokenmap["expires_at"] = utils.FormatTime(expires_at)
	td.SaveToken(tokenmap)
	return tokenmap, err
}

/**
 * Persistes auth tokens to file so it can be used later on.
 */
func (td *Client) SaveToken(token utils.StringMap) error {
	marshalled, err := json.Marshal(token)
	if err != nil {
		log.Printf("Could not marshall token: %+v", token)
		return err
	}
	data := []byte(marshalled)
	err = os.WriteFile(td.TokenFilePath(), data, 0777)
	return err
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

func (td *Client) GetChain(symbol string, date string, is_call bool, refresh_type int32) (*models.Chain, error) {
	if refresh_type == 0 {
		refresh_type = td.DefaultChainRefresh
	}
	date = strings.Replace(date, "-", "_", -1)
	qdb := td.db
	now := time.Now().UTC()
	chain := qdb.GetChain(symbol, date, is_call)
	var err error
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
	tickers := make(map[string]*models.Ticker)
	var tail []string
	for len(symbols) > 0 {
		numsyms := 200
		if len(symbols) < numsyms {
			numsyms = len(symbols)
		}
		symbols, tail = symbols[:numsyms], symbols[numsyms:]
		symstr := strings.Join(symbols, "%2C")
		url := fmt.Sprintf("%s?apikey=%s&symbol=%s", TDAMT_OPT_TICKER_URL, td.ClientId, symstr)

		result, _, err := utils.JsonGet(url, func(req *http.Request) {
			req.Header.Set("Authorization", td.Bearer())
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
	log.Printf("Loading chain data from server for %s...", symbol)
	url := fmt.Sprintf("%s?apikey=%s&symbol=%s", TDAMT_OPT_CHAIN_URL, td.ClientId, symbol)
	result, _, err := utils.JsonGet(url, func(req *http.Request) {
		req.Header.Set("Authorization", td.Bearer())
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

////////////////////////////////////////////////////////////////////////
//			Things related to streaming API and user principals
////////////////////////////////////////////////////////////////////////

func (td *Client) UserPrincipals() (utils.StringMap, error) {
	if td.userPrincipals == nil {
		if !td.IsAuthenticated() {
			return nil, fmt.Errorf("TD Client needs auth or tokens refreshed")
		}
		var err error
		if td.userPrincipals == nil {
			td.userPrincipals, err = td.FetchUserPrincipals()
			if err != nil || td.userPrincipals["error"] != nil {
				td.userPrincipals = nil
				log.Print("Error getting user principals: ", err, td.userPrincipals)
				return nil, err
			}
		}
	}
	return td.userPrincipals, nil
}

func (td *Client) FetchUserPrincipals() (utils.StringMap, error) {
	fstr := "streamerSubscriptionKeys,streamerConnectionInfo"
	url := fmt.Sprintf("%s?apikey=%s&fields=%s", TDAMT_USER_PRINCIPALS_URL, td.ClientId, fstr)
	result, _, err := utils.JsonGet(url, func(req *http.Request) {
		req.Header.Set("Authorization", td.Bearer())
	})
	if err != nil {
		log.Println("Error Fetching Principals: ", err)
	}
	return result.(utils.StringMap), err
}

/**
 * Gets the credentials required to connect to the streaming server.
 * This method also checks that a valid login/auth is available and
 * will try to refresh any tokens if required before failing on errors.
 */
func (td *Client) StreamingCredentials() (utils.StringMap, error) {
	if td.credentials == nil {
		userPrincipals, err := td.UserPrincipals()
		if err != nil {
			return nil, err
		}
		td.credentials, err = CredentialsFromPrincipal(userPrincipals)
		if err != nil {
			return nil, err
		}
		streamerInfo := td.userPrincipals["streamerInfo"].(utils.StringMap)
		td.wsUrl = &url.URL{
			Scheme: "wss",
			Host:   streamerInfo["streamerSocketUrl"].(string),
			Path:   "/ws",
		}
	}
	return td.credentials, nil
}

/**
 * Returns the WS connection URL.
 */
func (td *Client) WSUrl() (*url.URL, error) {
	if _, err := td.StreamingCredentials(); err != nil {
		return nil, err
	}
	return td.wsUrl, nil
}

func CredentialsFromPrincipal(userPrincipalsResponse utils.StringMap) (utils.StringMap, error) {
	layout := "2006-01-02T15:04:05+0000"
	streamerInfo := userPrincipalsResponse["streamerInfo"].(utils.StringMap)
	tokenTimeStampAsDateObj, err := time.Parse(layout, streamerInfo["tokenTimestamp"].(string))
	if err != nil {
		return nil, err
	}
	tokenTimeStampAsMs := tokenTimeStampAsDateObj.UnixMilli()
	credentials := make(map[string]interface{})
	account := userPrincipalsResponse["accounts"].([]interface{})[0].(utils.StringMap)
	credentials["authorized"] = "Y"
	credentials["timestamp"] = tokenTimeStampAsMs
	credentials["userid"] = account["accountId"].(string)
	credentials["company"] = account["company"].(string)
	credentials["segment"] = account["segment"].(string)
	credentials["cddomain"] = account["accountCdDomainId"].(string)
	credentials["usergroup"] = streamerInfo["userGroup"].(string)
	credentials["token"] = streamerInfo["token"].(string)
	credentials["accesslevel"] = streamerInfo["accessLevel"].(string)
	credentials["appid"] = streamerInfo["appId"].(string)
	credentials["acl"] = streamerInfo["acl"]
	return credentials, err
}
