package db

import (
	"gnuplex/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	ORM *gorm.DB
}

func Init(path string, verbose bool) (*DB, error) {
	logLevel := logger.Warn
	if verbose {
		logLevel = logger.Info
	}
	orm, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	db := DB{ORM: orm}
	if err != nil {
		return nil, err
	}
	// migrations
	if err := db.ORM.AutoMigrate(&models.MediaItem{}); err != nil {
		return nil, err
	}
	if err := db.ORM.AutoMigrate(&models.MediaDir{}); err != nil {
		return nil, err
	}
	if err := db.ORM.AutoMigrate(&models.FileExtension{}); err != nil {
		return nil, err
	}
	if err := db.ORM.AutoMigrate(&models.Settings{}); err != nil {
		return nil, err
	}
	return &db, nil
}
