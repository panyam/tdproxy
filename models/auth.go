package models

import (
	"encoding/json"
	"github.com/panyam/goutils/utils"
	"gorm.io/gorm"
	"time"
)

type Auth struct {
	ClientId           string `gorm:"primaryKey"`
	CallbackUrl        string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	ExpiresAt          time.Time
	AuthTokenJson      json.RawMessage
	UserPrincipalsJson json.RawMessage `gorm:"-"`
	authToken          utils.StringMap `gorm:"-"`
	userPrincipals     utils.StringMap `gorm:"-"`
}

func (a *Auth) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["client_id"] = a.ClientId
	out["callback_url"] = a.CallbackUrl
	out["auth_token"] = a.authToken
	out["user_principals"] = a.userPrincipals
	out["expires_at"] = utils.FormatTime(a.ExpiresAt)
	return out
}

func (auth *Auth) FromJson(json utils.StringMap) {
	if json != nil {
		auth.ClientId = json["client_id"].(string)
		auth.CallbackUrl = json["callback_url"].(string)
		if val, ok := json["auth_token"]; ok && val != nil {
			auth.SetAuthToken(val.(utils.StringMap))
		}
		if val, ok := json["user_principals"]; ok && val != nil {
			auth.SetUserPrincipals(val.(utils.StringMap))
		}
		if val, ok := json["expires_at"]; ok && val != nil {
			auth.ExpiresAt = utils.ParseTime(val.(string))
		}
	}
}

func (auth *Auth) AuthToken() utils.StringMap {
	return auth.authToken
}

func (auth *Auth) UserPrincipals() utils.StringMap {
	return auth.userPrincipals
}

func (auth *Auth) SetUserPrincipals(info utils.StringMap) bool {
	j, _ := json.Marshal(info)
	auth.userPrincipals = info
	auth.UserPrincipalsJson = j
	return true
}

func (auth *Auth) SetAuthToken(info utils.StringMap) bool {
	j, _ := json.Marshal(info)
	auth.authToken = info
	auth.AuthTokenJson = j
	return true
}

func (auth *Auth) IsAuthenticated() bool {
	if auth.authToken == nil {
		return false
	}
	if auth.ExpiresAt.Sub(time.Now().UTC()) <= 0 {
		return false
	}
	return true
}

func (auth *Auth) AccessToken() string {
	access_token := auth.authToken["access_token"]
	if access_token == nil {
		return ""
	}
	return access_token.(string)
}

func (auth *Auth) AfterFind(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	var res interface{}
	if auth.UserPrincipalsJson != nil {
		res, err = utils.JsonDecodeBytes(auth.UserPrincipalsJson)
		if err == nil && res != nil {
			auth.userPrincipals = res.(utils.StringMap)
		}
	}
	if auth.AuthTokenJson != nil {
		res, err = utils.JsonDecodeBytes(auth.AuthTokenJson)
		if err == nil && res != nil {
			auth.authToken = res.(utils.StringMap)
		}
	}
	return nil
}
