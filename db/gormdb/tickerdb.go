package gormdb

import (
	"errors"
	"gorm.io/gorm"
	"tdproxy/models"
)

type TickerDB struct {
	db *gorm.DB
}

func NewTickerDB(db *gorm.DB) *TickerDB {
	db.AutoMigrate(&models.Ticker{})
	return &TickerDB{db: db}
}

func (db *TickerDB) GetTicker(symbol string) (*models.Ticker, error) {
	var out models.Ticker
	err := db.db.First(&out, "symbol = ?", symbol).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &out, err
}

func (db *TickerDB) SaveTicker(ticker *models.Ticker) (err error) {
	result := db.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(ticker)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(ticker).Error
	}
	return
}
