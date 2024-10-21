package db

import (
	"gnuplex-backend/models"
	"log"
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

func (db *DB) SetFileExts(mediadirs []string) error {
	err := db.ORM.Unscoped().Where("1 = 1").Delete(&models.FileExtension{}).Error
	if err != nil {
		return err
	}
	for _, dir := range mediadirs {
		db.ORM.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.FileExtension{Extension: dir})
	}
	return err
}

func (db *DB) GetMediaDirs() ([]models.MediaDir, error) {
	var mediaDirs []models.MediaDir
	err := db.ORM.Order("path").Find(&mediaDirs).Error
	return mediaDirs, err
}

func (db *DB) SetMediadirs(mediaDirs []string) error {
	mediaDirsDB, err := db.GetMediaDirs()
	if err != nil {
		return err
	}
	mediaDirsH := make(map[string]bool)
	for _, dir := range mediaDirs {
		log.Println("a", dir)
		mediaDirsH[dir] = true
	}
	for _, dir := range mediaDirs {
		log.Println("b", dir)
		err := db.ORM.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.MediaDir{Path: dir}).Error
		if err != nil {
			log.Println("b", dir)
			return err
		}
	}
	for _, dir := range mediaDirsDB {
		if _, ok := mediaDirsH[dir.Path]; !ok {
			err := db.ORM.Unscoped().Delete(&dir).Error
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (db *DB) GetLast25Played() ([]models.MediaItem, error) {
	var mediaItems []models.MediaItem
	err := db.ORM.
		Order("last_played desc").
		Limit(25).
		Where("last_played != ''").
		Find(&mediaItems).Error
	return mediaItems, err
}

func (db *DB) GetMediaItems() ([]models.MediaItem, error) {
	var mediaItems []models.MediaItem
	err := db.ORM.
		Order("path").
		Find(&mediaItems).Error
	return mediaItems, err
}

func (db *DB) DeleteMediaItem(mediaItem *models.MediaItem) error {
	return db.ORM.Unscoped().Delete(mediaItem).Error
}

func (db *DB) DeleteMediaItemByPath(path string) error {
	return db.ORM.
		Unscoped().
		Where("path = ?", path).
		Delete(&models.MediaItem{}).Error
}

func (db *DB) AddMediaItem(path string) error {
	return db.ORM.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "path"}},
			DoNothing: true}).
		Create(&models.MediaItem{Path: path}).Error
}
