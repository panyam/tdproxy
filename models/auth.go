package models

import (
	"github.com/panyam/goutils/utils"
	"log"
	"time"
)

type AuthTokenJsonField struct {
	*Json
	AuthClientId string `gorm:"primaryKey"`
}

type UserPrincipalsJsonField struct {
	*Json
	AuthClientId string `gorm:"primaryKey"`
}

type Auth struct {
	ClientId              string `gorm:"primaryKey"`
	CallbackUrl           string
	ExpiresAt             time.Time
	RefreshTokenExpiresAt time.Time
	AuthToken             JsonField `gorm:"type:text"` // AuthTokenJsonField
	UserPrincipals        JsonField `gorm:"type:text"` // UserPrincipalsJsonField
}

func (a *Auth) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["client_id"] = a.ClientId
	out["callback_url"] = a.CallbackUrl
	out["auth_token"] = a.AuthTokenValue()
	out["user_principals"] = a.UserPrincipalsValue()
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

func (auth *Auth) AuthTokenValue() utils.StringMap {
	return auth.AuthToken
	/*
		res, err := auth.AuthToken.Value()
		if err != nil || res == nil {
			return nil
		}
		return res.(utils.StringMap)
	*/
}

func (auth *Auth) UserPrincipalsValue() utils.StringMap {
	return auth.UserPrincipals
	/*
		res, err := auth.UserPrincipals.Value()
		if err != nil || res == nil {
			return nil
		}
		return res.(utils.StringMap)
	*/
}

func (auth *Auth) SetUserPrincipals(info utils.StringMap) bool {
	// auth.userPrincipals = NewJson(fmt.Sprintf("auth_%s_up", auth.ClientId), info)
	auth.UserPrincipals = JsonField(info)
	/*
		UserPrincipalsJsonField{
			AuthClientId: auth.ClientId,
			Json:         NewJson(info),
		}
	*/
	return true
}

func (auth *Auth) SetAuthToken(info utils.StringMap) bool {
	auth.AuthToken = JsonField(info)
	/*AuthTokenJsonField{
		AuthClientId: auth.ClientId,
		Json:         NewJson(info),
	}
	*/

	now := time.Now().UTC()
	expires_in := time.Duration(0)
	if val, ok := info["expires_in"]; ok {
		expires_in = time.Duration(val.(float64))
	} else {
		log.Println("Could not find expires_in: ", info)
	}
	auth.ExpiresAt = now.Add(expires_in * time.Second)
	if val, ok := info["refresh_token_expires_in"]; ok && val != nil {
		refresh_token_expires_in := time.Duration(val.(float64))
		auth.RefreshTokenExpiresAt = now.Add(refresh_token_expires_in * time.Second)
	}
	log.Println("Now, ExpiresIn, ExpiresAt: ",
		now, expires_in, auth.ExpiresAt,
		auth.RefreshTokenExpiresAt)
	return true
}

func (auth *Auth) IsAuthenticated() bool {
	if auth.AuthToken == nil {
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
	access_token := auth.AuthTokenValue()["access_token"]
	if access_token == nil {
		return ""
	}
	return access_token.(string)
}
