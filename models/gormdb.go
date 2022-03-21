package models

import (
	"fmt"
	"gorm.io/gorm"
)

type JsonValue struct {
	Key       string `gorm:"primaryKey"`
	ValueJson string
}

func (ticker *Ticker) AfterFind(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	// Do nothing here - we just create a Json entry with the right key
	ticker.Info = &Json{Key: ticker.Symbol, db: tx}
	return nil
}

func (option *Option) AfterFind(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	option.Info = &Json{Key: option.ShortKey(), db: tx}
	return
}

func (auth *Auth) AfterFind(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	auth.authToken = &Json{Key: fmt.Sprintf("auth_%s_at", auth.ClientId), db: tx}
	auth.userPrincipals = &Json{Key: fmt.Sprintf("auth_%s_up", auth.ClientId), db: tx}
	return nil
}

func (ticker *Ticker) AfterSave(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	err = ticker.Info.Save(tx)
	return
}

func (option *Option) AfterSave(tx *gorm.DB) (err error) {
	// Updated Stuff from json fields
	err = option.Info.Save(tx)
	return
}

func (auth *Auth) AfterSave(tx *gorm.DB) (err error) {
	// TODO - save transactionally
	err = auth.authToken.Save(tx)
	if err == nil {
		err = auth.userPrincipals.Save(tx)
	}
	return
}
