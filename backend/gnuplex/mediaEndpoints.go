package gnuplex

import (
	"errors"
	"gnuplex/models"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (gnuplex *GNUPlex) ScanLib() error {
	lastScanUUID := uuid.New().String()
	log.Println("scanuuid", lastScanUUID)
	/*
	 * grab MediaDirs, MediaItems, FileExts from the database
	 */
	mediaDirs, err := gnuplex.DB.GetMediaDirs()
	if err != nil {
		return err
	}
	fileExts, err := gnuplex.DB.GetFileExts()
	if err != nil {
		return err
	}
	fileExtH := make(map[string]bool)
	for _, fileExt := range fileExts {
		fileExtH[fileExt.Extension] = true
	}
	/*
	 * add new stuff
	 */
	for _, mediaDir := range mediaDirs {
		dir, err := os.Stat(mediaDir.Path)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediaDir.Path, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					log.Println("Walkdir prob: ", mediaDir)
					return err
				} else if !entry.IsDir() {
					ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
					if _, match := fileExtH[ext]; !match {
						return gnuplex.DB.AddMediaItemFile(path, lastScanUUID)
					}
					return nil
				} else {
					return nil
				}
			})
			if err != nil {
				return err
			}
		}
	}
	/*
	 * remove stuff that no longer exists
	 */
	return gnuplex.DB.DeleteMediaItemFilesNotMatchingUUID(lastScanUUID)
}

func (gnuplex *GNUPlex) NowPlaying() (*models.MediaItem, error) {
	if len(gnuplex.PlayQueue) == 0 {
		return nil, errors.New("PlayQueue is empty at the moment")
	}
	return gnuplex.PlayQueue[0], nil
}

func (gnuplex *GNUPlex) ReplaceQueueAndPlay(id models.MediaItemId) error {
	var mediaItem *models.MediaItem
	if err := gnuplex.DB.ORM.First(&mediaItem, id).Error; err != nil {
		return err
	}
	if mediaItem != nil {
		gnuplex.PlayQueue = []*models.MediaItem{mediaItem}
	}
	if err := gnuplex.MPV.ReplaceQueueAndPlay(mediaItem.Path); err != nil {
		return err
	}
	if err := gnuplex.DB.UpdateLastPlayed(mediaItem); err != nil {
		return err
	}
	return nil
}

func (gnuplex *GNUPlex) ReplaceQueueAndPlayByPath(path string) error {
	var mediaItem *models.MediaItem
	if err := gnuplex.DB.ORM.First(&mediaItem, "path = ?", path).Error; err != nil {
		return err
	}
	if mediaItem != nil {
		gnuplex.PlayQueue = []*models.MediaItem{mediaItem}
	}
	if err := gnuplex.MPV.ReplaceQueueAndPlay(mediaItem.Path); err != nil {
		return err
	}
	if err := gnuplex.DB.UpdateLastPlayed(mediaItem); err != nil {
		return err
	}
	return nil
}

func (gnuplex *GNUPlex) ReplaceQueueAndCastTempUrl(url string) error {
	return gnuplex.MPV.ReplaceQueueAndPlay(url)
}

func (gnuplex *GNUPlex) QueueLast(id models.MediaItemId) *models.MediaItem {
	var mediaItem *models.MediaItem
	gnuplex.DB.ORM.First(&mediaItem, id)
	if mediaItem != nil {
		gnuplex.PlayQueue = append(gnuplex.PlayQueue, mediaItem)
	}
	gnuplex.MPV.QueueMedia(mediaItem.Path)
	return mediaItem
}

func (gnuplex *GNUPlex) Cast(url string, temp bool) error {
	if temp {
		return gnuplex.ReplaceQueueAndCastTempUrl(url)
	} else if err := gnuplex.DB.AddMediaItemURL(url); err != nil {
		return err
	} else {
		return gnuplex.ReplaceQueueAndPlayByPath(url)
	}
}
