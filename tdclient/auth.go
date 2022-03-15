package tdclient

import (
	"encoding/json"
	"fmt"
	"github.com/panyam/goutils/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type AuthStore struct {
	RootDir  string
	auths    map[string]*Auth
	lastAuth *Auth
}

/**
 * Creates a new auth store.
 */
func NewAuthStore(rootdir string) *AuthStore {
	fmt.Println("Client Root Dir: ", rootdir)
	out := &AuthStore{RootDir: rootdir}
	os.MkdirAll(rootdir, 0777)
	out.Reload()
	return out
}

/**
 * Checks if a particular client id is authenticated.
 */
func (a *AuthStore) IsAuthenticated(client_id string) bool {
	if _, err := a.Reload(); err != nil {
		log.Println("Error loading auth tokens")
		return false
	}
	if auth, ok := a.auths[client_id]; ok && auth != nil {
		return auth.IsAuthenticated()
	}
	return false
}

func (a *AuthStore) TokensFilePath() string {
	return path.Join(a.RootDir, "tokens")
}

/**
 * Persistes auth tokens to file so it can be used later on.
 */
func (a *AuthStore) SaveTokens() (err error) {
	log.Println("Saving tokens...")
	defer log.Println("Finished Saved Tokens, err: ", err)
	auths := make(utils.StringMap)
	for key, value := range a.auths {
		auths[key] = value.ToJson()
	}
	var marshalled []byte
	marshalled, err = json.MarshalIndent(auths, "", "  ")
	if err != nil {
		log.Printf("Could not marshall token: %+v", a.auths)
		return err
	}
	err = os.WriteFile(a.TokensFilePath(), marshalled, 0777)
	return err
}

/**
 * Reloads the auth store contents.
 */
func (a *AuthStore) Reload() (auths map[string]*Auth, err error) {
	auths = make(map[string]*Auth)
	contents, err := os.ReadFile(a.TokensFilePath())
	if err != nil {
		log.Println(err)
		return auths, err
	}
	tokens, err := utils.JsonDecodeBytes(contents)
	if err != nil {
		log.Println(err)
		return auths, err
	}
	entries := tokens.(utils.StringMap)
	for clientId, entry := range entries {
		clientInfo := entry.(utils.StringMap)
		callback_url := clientInfo["callback_url"]
		auth := a.EnsureAuth(clientId, callback_url.(string))
		auth.FromJson(clientInfo)
	}
	fmt.Println("Loaded auth tokens: ", entries)
	return auths, nil
}

/**
 * Creates a new auth object and adds to the store.
 */
func (a *AuthStore) EnsureAuth(client_id string, callback_url string) (auth *Auth) {
	auth, ok := a.auths[client_id]
	if !ok || auth == nil {
		auth = &Auth{ClientId: client_id, CallbackUrl: callback_url}
		if a.auths == nil {
			a.auths = make(map[string]*Auth)
		}
		a.auths[client_id] = auth
	}
	auth.ClientId = client_id
	auth.CallbackUrl = callback_url
	a.lastAuth = auth
	return
}

func (a *AuthStore) LastAuth() *Auth {
	return a.lastAuth
}

const (
	TDAMT_AUTH_URL            = "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=%s&client_id=%s%%40AMER.OAUTHAP"
	TDAMT_TOKEN_URL           = "https://api.tdameritrade.com/v1/oauth2/token"
	TDAMT_USER_PRINCIPALS_URL = "https://api.tdameritrade.com/v1/userprincipals"
)

type Auth struct {
	ClientId       string
	CallbackUrl    string
	authToken      utils.StringMap
	userPrincipals utils.StringMap
	credentials    utils.StringMap
	wsUrl          *url.URL
}

func (a *Auth) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["client_id"] = a.ClientId
	out["callback_url"] = a.CallbackUrl
	out["auth_token"] = a.authToken
	out["user_principals"] = a.userPrincipals
	return out
}

func (auth *Auth) FromJson(json utils.StringMap) {
	if json != nil {
		auth.ClientId = json["client_id"].(string)
		auth.CallbackUrl = json["callback_url"].(string)
		if val, ok := json["auth_token"]; ok && val != nil {
			auth.authToken = val.(utils.StringMap)
		}
		if val, ok := json["user_principals"]; ok && val != nil {
			auth.userPrincipals = val.(utils.StringMap)
		}
	}
}

func (auth *Auth) Bearer() string {
	return fmt.Sprintf("Bearer %s", auth.GetAccessToken())
}

func (auth *Auth) GetAccessToken() string {
	access_token := auth.authToken["access_token"]
	if access_token == nil {
		return ""
	}
	return access_token.(string)
}

func (auth *Auth) StartAuthUrl() string {
	callback_quoted := url.QueryEscape(auth.CallbackUrl)
	url := fmt.Sprintf(TDAMT_AUTH_URL, callback_quoted, auth.ClientId)
	return url
}

func (auth *Auth) CompleteAuth(code string) (err error) {
	log.Println("Completing Auth...")
	defer func() {
		log.Println("Completed auth, err: ", err)
		if err != nil {
			auth.authToken = nil
			auth.userPrincipals = nil
		}
	}()
	// decoded := code
	var decoded string
	decoded, err = url.PathUnescape(code)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", auth.ClientId)
	form.Add("redirect_uri", auth.CallbackUrl)
	form.Add("code", decoded)
	fmt.Println("Form: ", form)
	var postreq *http.Request
	postreq, err = http.NewRequest("POST", TDAMT_TOKEN_URL, strings.NewReader(form.Encode()))
	postreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	var response *http.Response
	response, err = client.Do(postreq)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	if response.StatusCode != 200 {
		fmt.Println("Failed Response: ", response)
		return fmt.Errorf(response.Status)
	}

	// Save it locally now
	var token interface{}
	token, err = utils.JsonDecodeBytes(body)
	if err != nil {
		fmt.Println("Invalid response json: ", err)
		return err
	}
	tokenmap := token.(utils.StringMap)
	expires_in := time.Duration(tokenmap["expires_in"].(float64))
	expires_at := now.Add(expires_in * time.Second)
	log.Println("Now, ExpiresIn, ExpiresAt: ", now, expires_in, expires_at)
	tokenmap["expires_at"] = utils.FormatTime(expires_at)
	auth.authToken = tokenmap
	return err
}

////////////////////////////////////////////////////////////////////////
//			Things related to streaming API and user principals
////////////////////////////////////////////////////////////////////////

func (auth *Auth) UserPrincipals() (utils.StringMap, error) {
	if auth.userPrincipals == nil {
		if !auth.IsAuthenticated() {
			return nil, fmt.Errorf("TD Client needs auth or tokens refreshed")
		}
		var err error
		if auth.userPrincipals == nil {
			auth.userPrincipals, err = auth.FetchUserPrincipals()
			if err != nil || auth.userPrincipals["error"] != nil {
				auth.userPrincipals = nil
				log.Print("Error getting user principals: ", err, auth.userPrincipals)
				return nil, err
			}
		}
	}
	return auth.userPrincipals, nil
}

func (auth *Auth) IsAuthenticated() bool {
	if auth.authToken == nil {
		return false
	}
	expires_at_str := auth.authToken["expires_at"]
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

func (auth *Auth) FetchUserPrincipals() (utils.StringMap, error) {
	fstr := "streamerSubscriptionKeys,streamerConnectionInfo"
	url := fmt.Sprintf("%s?apikey=%s&fields=%s", TDAMT_USER_PRINCIPALS_URL, auth.ClientId, fstr)
	result, _, err := utils.JsonGet(url, func(req *http.Request) {
		req.Header.Set("Authorization", auth.Bearer())
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
func (auth *Auth) StreamingCredentials() (utils.StringMap, error) {
	if auth.credentials == nil {
		userPrincipals, err := auth.UserPrincipals()
		if err != nil {
			return nil, err
		}
		auth.credentials, err = CredentialsFromPrincipal(userPrincipals)
		if err != nil {
			return nil, err
		}
		streamerInfo := auth.userPrincipals["streamerInfo"].(utils.StringMap)
		auth.wsUrl = &url.URL{
			Scheme: "wss",
			Host:   streamerInfo["streamerSocketUrl"].(string),
			Path:   "/ws",
		}
	}
	return auth.credentials, nil
}

/**
 * Returns the WS connection URL.
 */
func (auth *Auth) WSUrl() (*url.URL, error) {
	if _, err := auth.StreamingCredentials(); err != nil {
		return nil, err
	}
	return auth.wsUrl, nil
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
