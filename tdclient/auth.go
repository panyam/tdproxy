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
func NewAuthStore(authdb db.AuthDB) (a *AuthStore) {
	a = &AuthStore{authdb: authdb}
	last := authdb.LastAuth()
	if last != nil {
		a.lastAuth = &Auth{Auth: last}
	}
	return
}

/**
 * Checks if a particular client id is authenticated.
 */
func (a *AuthStore) EnsureAuthenticated(client_id string) bool {
	auth, err := a.EnsureAuth(client_id, "")
	if auth == nil || err != nil {
		return false
	}
	res := auth.IsAuthenticated()
	if res {
		a.lastAuth = auth
	} else if auth.CanRefreshToken() {
		if err := auth.RefreshTokens(); err == nil {
			res = true
			a.lastAuth = auth
			a.SaveAuth(auth)
		}
	}
	return res
}

func (a *AuthStore) EnsureAuth(client_id string, callback_url string) (*Auth, error) {
	auth, err := a.authdb.EnsureAuth(client_id)
	if err != nil {
		return nil, err
	}
	if callback_url != "" {
		auth.CallbackUrl = callback_url
	}
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

func (auth *Auth) RefreshTokens() (err error) {
	log.Println("Refreshing Tokens...")
	defer func() {
		log.Println("Completed auth, err: ", err)
		if err != nil {
			auth.SetUserPrincipals(nil)
		}
	}()
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", auth.AuthTokenValue()["refresh_token"].(string))
	form.Add("client_id", auth.ClientId)
	form.Add("redirect_uri", auth.CallbackUrl)
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

	fmt.Println("Refresh Tokens Response Status:", response.Status)
	fmt.Println("Refresh Tokens Response Headers:", response.Header)
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
	auth.SetAuthToken(token.(utils.StringMap))
	return err
}

func (auth *Auth) CompleteAuth(code string) (err error) {
	log.Println("Completing Auth...")
	defer func() {
		log.Println("Completed auth, err: ", err)
		if err != nil {
			auth.SetAuthToken(nil)
			auth.SetUserPrincipals(nil)
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
	form.Add("access_type", "offline")
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
	auth.SetAuthToken(token.(utils.StringMap))
	return err
}

////////////////////////////////////////////////////////////////////////
//			Things related to streaming API and user principals
////////////////////////////////////////////////////////////////////////

func (auth *Auth) EnsureUserPrincipals() error {
	if auth.UserPrincipalsValue() == nil {
		if !auth.IsAuthenticated() {
			return fmt.Errorf("TD Client needs auth or tokens refreshed")
		}
		if auth.UserPrincipalsValue() == nil {
			up, err := auth.FetchUserPrincipals()
			if err != nil || up["error"] != nil {
				auth.SetUserPrincipals(nil)
				log.Print("Error getting user principals: ", err, auth.UserPrincipalsValue())
				return err
			} else {
				auth.SetUserPrincipals(up)
			}
		}
	}
	return nil
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
func (auth *Auth) StreamingCredentials() (creds utils.StringMap, err error) {
	if auth.credentials == nil {
		if err = auth.EnsureUserPrincipals(); err != nil {
			return nil, err
		}
		up := auth.UserPrincipalsValue()
		auth.credentials, err = CredentialsFromPrincipal(up)
		if err != nil {
			return nil, err
		}
		streamerInfo := up["streamerInfo"].(utils.StringMap)
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
