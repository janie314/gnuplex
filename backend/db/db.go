package db

import (
	"gnuplex-backend/models"
	"time"

	"gorm.io/gorm"
)

func Init(db *gorm.DB) error {
	// migrations
	err := db.AutoMigrate(&models.MediaItem{})
	return err
}

func UpdateLastPlayed(db *gorm.DB, mediaItem *models.MediaItem) {
	db.Model(mediaItem).Update("LastPlayed", time.Now().UTC().Format(time.RFC3339))
}
