package models

import "gorm.io/gorm"

type MediaItem struct {
	gorm.Model
	Path       string
	LastPlayed string `gorm:"index:,sort:desc"`
}

type MediaItemId uint
