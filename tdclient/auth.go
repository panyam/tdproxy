package tdclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"legfinder/tdproxy/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Auth struct {
	RootDir        string
	ClientId       string
	CallbackUrl    string
	authToken      utils.StringMap
	userPrincipals utils.StringMap
	credentials    utils.StringMap
	wsUrl          *url.URL
}

const (
	TDAMT_AUTH_URL            = "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=%s&client_id=%s%%40AMER.OAUTHAP"
	TDAMT_TOKEN_URL           = "https://api.tdameritrade.com/v1/oauth2/token"
	TDAMT_USER_PRINCIPALS_URL = "https://api.tdameritrade.com/v1/userprincipals"
)

func NewAuth(rootdir string, client_id string, callback_url string) *Auth {
	fmt.Println("Client Root Dir: ", rootdir)
	out := &Auth{RootDir: rootdir, ClientId: client_id, CallbackUrl: callback_url}
	os.MkdirAll(rootdir, 0777)
	out.ReloadAuthTokens()
	return out
}

func (td *Auth) IsAuthenticated() bool {
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

func (td *Auth) TokenFilePath() string {
	return path.Join(td.RootDir, "tokens")
}

func (td *Auth) Bearer() string {
	return fmt.Sprintf("Bearer %s", td.GetAccessToken())
}

func (td *Auth) GetAccessToken() string {
	return td.authToken["access_token"].(string)
}

func (td *Auth) ReloadAuthTokens() utils.StringMap {
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

func (td *Auth) StartAuthUrl() string {
	callback_quoted := url.QueryEscape(td.CallbackUrl)
	url := fmt.Sprintf(TDAMT_AUTH_URL, callback_quoted, td.ClientId)
	return url
}

func (td *Auth) CompleteAuth(code string) (utils.StringMap, error) {
	// decoded := code
	decoded, err := url.PathUnescape(code)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", td.ClientId)
	form.Add("redirect_uri", td.CallbackUrl)
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
func (td *Auth) SaveToken(token utils.StringMap) error {
	marshalled, err := json.Marshal(token)
	if err != nil {
		log.Printf("Could not marshall token: %+v", token)
		return err
	}
	data := []byte(marshalled)
	err = os.WriteFile(td.TokenFilePath(), data, 0777)
	return err
}

////////////////////////////////////////////////////////////////////////
//			Things related to streaming API and user principals
////////////////////////////////////////////////////////////////////////

func (td *Auth) UserPrincipals() (utils.StringMap, error) {
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

func (td *Auth) FetchUserPrincipals() (utils.StringMap, error) {
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
func (td *Auth) StreamingCredentials() (utils.StringMap, error) {
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
func (td *Auth) WSUrl() (*url.URL, error) {
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
