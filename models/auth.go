package models

import (
	"encoding/json"
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
	"time"
)

type Auth struct {
	ClientId              string `gorm:"primaryKey"`
	CallbackUrl           string
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
	ExpiresAt             time.Time
	RefreshTokenExpiresAt time.Time
	AuthTokenJson         json.RawMessage
	UserPrincipalsJson    json.RawMessage `gorm:"-"`
	authToken             *Json
	userPrincipals        *Json
}

func (a *Auth) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["client_id"] = a.ClientId
	out["callback_url"] = a.CallbackUrl
	out["auth_token"] = a.authToken
	out["user_principals"] = a.userPrincipals
	out["expires_at"] = utils.FormatTime(a.ExpiresAt)
	out["refresh_token_expires_at"] = utils.FormatTime(a.ExpiresAt)
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
		if val, ok := json["refresh_token_expires_at"]; ok && val != nil {
			auth.RefreshTokenExpiresAt = utils.ParseTime(val.(string))
		}
	}
}

func (auth *Auth) AuthToken() utils.StringMap {
	res, err := auth.authToken.Value()
	if err != nil || res == nil {
		return nil
	}
	return res.(utils.StringMap)
}

func (auth *Auth) UserPrincipals() utils.StringMap {
	res, err := auth.userPrincipals.Value()
	if err != nil || res == nil {
		return nil
	}
	return res.(utils.StringMap)
}

func (auth *Auth) SetUserPrincipals(info utils.StringMap) bool {
	auth.userPrincipals = NewJson(fmt.Sprintf("auth_%s_up", auth.ClientId), info)
	return true
}

func (auth *Auth) SetAuthToken(info utils.StringMap) bool {
	auth.authToken = NewJson(fmt.Sprintf("auth_%s_at", auth.ClientId), info)

	now := time.Now().UTC()
	expires_in := time.Duration(info["expires_in"].(float64))
	refresh_token_expires_in := time.Duration(info["refresh_token_expires_in"].(float64))
	auth.ExpiresAt = now.Add(expires_in * time.Second)
	auth.RefreshTokenExpiresAt = now.Add(refresh_token_expires_in * time.Second)
	log.Println("Now, ExpiresIn, ExpiresAt: ", now, expires_in, auth.ExpiresAt)
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

/**
 * Check if refresh token is valid.
 */
func (auth *Auth) CanRefreshToken() bool {
	if auth.RefreshTokenExpiresAt.Sub(time.Now().UTC()) <= 0 {
		return false
	}
	return true
}

func (auth *Auth) AccessToken() string {
	access_token := auth.AuthToken()["access_token"]
	if access_token == nil {
		return ""
	}
	return access_token.(string)
}
