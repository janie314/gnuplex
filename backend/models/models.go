package models

import "gorm.io/gorm"

type MediaItem struct {
	gorm.Model
	Path       string `gorm:"index"`
	LastPlayed string `gorm:"index:,sort:desc"`
}

type MediaItemId uint

type MediaDir struct {
	gorm.Model
	Path        string `gorm:"index"`
	LastScanned string `gorm:"index:,sort:desc"`
}

type MediaDirId uint

type FileExtension struct {
	gorm.Model
	Extension string
}

type FileExtensionId uint
