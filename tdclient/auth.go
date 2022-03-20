package tdclient

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"tdproxy/db"
	"tdproxy/models"
	"time"
)

type Auth struct {
	*models.Auth
	wsUrl       *url.URL
	credentials utils.StringMap
}

type AuthStore struct {
	authdb   db.AuthDB
	lastAuth *Auth
}

/**
 * Creates a new auth store.
 */
func NewAuthStore(authdb db.AuthDB) *AuthStore {
	return &AuthStore{authdb: authdb}
}

/**
 * Checks if a particular client id is authenticated.
 */
func (a *AuthStore) IsAuthenticated(client_id string) bool {
	auth, err := a.authdb.EnsureAuth(client_id)
	if auth == nil || err != nil {
		return false
	}
	res := auth.IsAuthenticated()
	if res {
		a.lastAuth = &Auth{Auth: auth}
	}
	return res
}

func (a *AuthStore) EnsureAuth(client_id string, callback_url string) (*Auth, error) {
	auth, err := a.authdb.EnsureAuth(client_id)
	if err != nil {
		return nil, err
	}
	auth.CallbackUrl = callback_url
	return &Auth{
		Auth: auth,
	}, nil
}

func (a *AuthStore) LastAuth() *Auth {
	return a.lastAuth
}

const (
	TDAMT_AUTH_URL            = "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=%s&client_id=%s%%40AMER.OAUTHAP"
	TDAMT_TOKEN_URL           = "https://api.tdameritrade.com/v1/oauth2/token"
	TDAMT_USER_PRINCIPALS_URL = "https://api.tdameritrade.com/v1/userprincipals"
)

func (auth *Auth) Bearer() string {
	return fmt.Sprintf("Bearer %s", auth.AccessToken())
}

func (a *AuthStore) SaveAuth(auth *Auth) (err error) {
	return a.authdb.SaveAuth(auth.Auth)
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
			auth.AuthToken = nil
			auth.UserPrincipals = nil
		}
	}()
	// decoded := code
	var decoded string
	decoded, err = url.PathUnescape(code)
	if err != nil {
		return err
	}
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
	now := time.Now().UTC()
	tokenmap := token.(utils.StringMap)
	expires_in := time.Duration(tokenmap["expires_in"].(float64))
	auth.ExpiresAt = now.Add(expires_in * time.Second)
	log.Println("Now, ExpiresIn, ExpiresAt: ", now, expires_in, auth.ExpiresAt)
	auth.AuthToken = tokenmap
	return err
}

////////////////////////////////////////////////////////////////////////
//			Things related to streaming API and user principals
////////////////////////////////////////////////////////////////////////

func (auth *Auth) EnsureUserPrincipals() (utils.StringMap, error) {
	if auth.UserPrincipals == nil {
		if !auth.IsAuthenticated() {
			return nil, fmt.Errorf("TD Client needs auth or tokens refreshed")
		}
		var err error
		if auth.UserPrincipals == nil {
			auth.UserPrincipals, err = auth.FetchUserPrincipals()
			if err != nil || auth.UserPrincipals["error"] != nil {
				auth.UserPrincipals = nil
				log.Print("Error getting user principals: ", err, auth.UserPrincipals)
				return nil, err
			}
		}
	}
	return auth.UserPrincipals, nil
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
		userPrincipals, err := auth.EnsureUserPrincipals()
		if err != nil {
			return nil, err
		}
		auth.credentials, err = CredentialsFromPrincipal(userPrincipals)
		if err != nil {
			return nil, err
		}
		streamerInfo := auth.UserPrincipals["streamerInfo"].(utils.StringMap)
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
