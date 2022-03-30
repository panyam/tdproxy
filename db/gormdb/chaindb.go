package gormdb

import (
	"log"
	// "gorm.io/driver/sqlite"
	"errors"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"tdproxy/models"
	"time"
)

type ChainDB struct {
	optiondb *OptionDB
	db       *gorm.DB
}

func NewChainDB(db *gorm.DB) *ChainDB {
	db.AutoMigrate(&models.ChainInfo{})
	db.AutoMigrate(&models.Chain{})
	return &ChainDB{
		optiondb: NewOptionDB(db),
		db:       db,
	}
}

func (db *ChainDB) GetChainInfo(symbol string) (*models.ChainInfo, error) {
	var out models.ChainInfo
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

func (db *ChainDB) SaveChainInfo(symbol string, last_refreshed_at time.Time) error {
	return nil
	//return db.db.Save(option).Error
}

func (db *ChainDB) GetChain(symbol string, date string, is_call bool) (*models.Chain, error) {
	var out models.Chain
	err := db.db.First(&out, "is_call = ? AND symbol = ? AND date_string = ?", is_call, symbol, date).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// now load options
	out.Options, err = db.optiondb.GetOptions(is_call, symbol, date)
	return &out, err
}

func (db *ChainDB) SaveChain(chain *models.Chain) error {
	result := db.db.Model(chain).Updates(chain)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(chain).Error
	}
	err = db.optiondb.SaveOptions(chain.Options)
	if err != nil {
		log.Println("Error Saving Options: ", err)
	}
	return err
}
