package models

import (
	"encoding/json"
	"errors"
	"github.com/panyam/goutils/utils"
	"gorm.io/gorm"
	"log"
)

type Json struct {
	Key       string
	ValueJson string // json.RawMessage
	value     interface{}
	loaded    bool
	db        *gorm.DB
}

func NewJson(key string, value interface{}) *Json {
	j, _ := json.Marshal(value)
	return &Json{ValueJson: string(j), value: value, loaded: true}
}

func (jr *Json) Value() (result interface{}, err error) {
	if jr.value == nil {
		err = jr.db.First(jr, "key = ?", jr.Key).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			} else {
				log.Println("Error finding jr: ", jr.Key, err)
				return nil, err
			}
		}
		jr.value, err = utils.JsonDecodeBytes([]byte(jr.ValueJson))
	}
	return jr.value, err
}

func (jr *Json) Save(db *gorm.DB) (err error) {
	result := db.Model(jr).Updates(jr)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.Create(jr).Error
	}
	return
}
