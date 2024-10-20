package models

import "gorm.io/gorm"

type MediaItem struct {
	gorm.Model
	Path       string
	LastPlayed string
}

type MediaItemId uint
