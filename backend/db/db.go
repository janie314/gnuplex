package db

import (
	"gnuplex-backend/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DB struct {
	ORM *gorm.DB
}

func Init(path string) (*DB, error) {
	orm, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
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
	return &db, nil
}

func (db *DB) UpdateLastPlayed(mediaItem *models.MediaItem) error {
	return db.ORM.Model(mediaItem).Update("LastPlayed", time.Now().UTC().Format(time.RFC3339)).Error
}

func (db *DB) GetFileExts() ([]models.FileExtension, error) {
	var exts []models.FileExtension
	err := db.ORM.Find(&exts).Error
	return exts, err
}

func (db *DB) GetMediaDirs() ([]models.MediaDir, error) {
	var mediaDirs []models.MediaDir
	err := db.ORM.Order("path").Find(&mediaDirs).Error
	return mediaDirs, err
}

func (db *DB) SetMediadirs(mediadirs []string) error {
	err := db.ORM.Where("1 = 1").Delete(&models.MediaDir{}).Error
	if err != nil {
		return err
	}
	for _, dir := range mediadirs {
		db.ORM.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.MediaDir{Path: dir})
	}
	return err
}

func (db *DB) GetMediaItems() ([]models.MediaItem, error) {
	var mediaItems []models.MediaItem
	err := db.ORM.Order("path").Find(&mediaItems).Error
	return mediaItems, err
}

func (db *DB) DeleteMediaItem(mediaItem *models.MediaItem) error {
	return db.ORM.Delete(mediaItem).Error
}

func (db *DB) AddMediaItem(path string) error {
	return db.ORM.Create(&models.MediaItem{Path: path}).Error
}
