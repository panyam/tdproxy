package models

import (
	"encoding/json"
	"errors"
	"github.com/panyam/goutils/utils"
)

var InvalidJsonKeyError = errors.New("InvalidJsonKeyError")

type Json struct {
	ValueJson     string
	LastUpdatedAt int64 `gorm:"autoUpdateTime:milli"`
	value         interface{}
}

func NewJson(value interface{}) *Json {
	j, _ := json.Marshal(value)
	return &Json{ValueJson: string(j), value: value}
}

func (jr *Json) HasValue() bool {
	if jr == nil {
		return false
	}
	var err error
	if jr.value == nil && jr.ValueJson != "" {
		jr.value, err = utils.JsonDecodeStr(jr.ValueJson)
	}
	return err == nil && jr.value != nil
}

func (jr *Json) Value() (result interface{}, err error) {
	if jr == nil {
		return nil, nil
	}
	if jr.value == nil && jr.ValueJson != "" {
		jr.value, err = utils.JsonDecodeStr(jr.ValueJson)
	}
	return jr.value, err
}
