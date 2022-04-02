package gormdb

import (
	"errors"
	"gorm.io/gorm"
	"tdproxy/models"
)

type OptionDB struct {
	db *gorm.DB
}

func NewOptionDB(db *gorm.DB) *OptionDB {
	db.AutoMigrate(&models.Option{})
	// db.AutoMigrate(&models.OptionJsonField{})
	return &OptionDB{db: db}
}

func (db *OptionDB) GetOption(symbol string, date string, is_call bool, price string) (*models.Option, error) {
	var option models.Option
	err := db.db.First(&option, "symbol = ? AND date_string = ? AND is_call = ? AND price_string = ?", symbol, date, is_call, price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &option, err
}

/**
 * Get options for a particular symbol on a given date for either calls or puts.
 */
func (db *OptionDB) GetOptions(is_call bool, symbol string, date string) ([]*models.Option, error) {
	var options []*models.Option
	err := db.db.Where("symbol = ? AND date_string = ? AND is_call = ?", symbol, date, is_call).Find(&options).Error
	if err != nil {
		return nil, err
	}
	return options, err
}

/**
 * Save a particular option.
 */
func (db *OptionDB) SaveOption(option *models.Option) (err error) {
	result := db.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(option)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(option).Error
	}
	return
}

/**
 * Delete options on a particular date for a given symbol
 */
func (db *OptionDB) DeleteOptions(symbol string, date string, is_call bool) error {
	return db.db.Where("symbol = ? AND date_string = ? AND is_call = ?", symbol, date, is_call).Delete(&models.Option{}).Error
}

func (db *OptionDB) SaveOptions(options []*models.Option) error {
	return db.db.Transaction(func(tx *gorm.DB) error {
		for _, option := range options {
			if err := db.SaveOption(option); err != nil {
				return err
			}
		}
		return nil
	})
}
