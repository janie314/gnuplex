package models

import "gorm.io/gorm"

type MediaItemType int

const (
	File MediaItemType = 1
	URL  MediaItemType = 2
)

type MediaItem struct {
	gorm.Model
	Path         string `gorm:"uniqueIndex"`
	LastPlayed   string `gorm:"index:,sort:desc"`
	Temp         bool   `gorm:"default:false"`
	Type         MediaItemType
	LastScanUUID string `gorm:"index:,default:''"`
}

type MediaItemId uint

type MediaDir struct {
	gorm.Model
	Path        string `gorm:"uniqueIndex"`
	LastScanned string `gorm:"index:,sort:desc"`
}

type MediaDirId uint

type FileExtension struct {
	gorm.Model
	Extension string
}

type FileExtensionId uint
